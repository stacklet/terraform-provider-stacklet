package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ datasource.DataSource = &ssoGroupDataSource{}
)

func NewSSOGroupDataSource() datasource.DataSource {
	return &ssoGroupDataSource{}
}

type ssoGroupDataSource struct {
	client *graphql.Client
}

type ssoGroupDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Roles             types.List   `tfsdk:"roles"`
	AccountGroupUUIDs types.List   `tfsdk:"account_group_uuids"`
}

func (d *ssoGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_group"
}

func (d *ssoGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an SSO group configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for this SSO group configuration.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name that identifies the group in the external SSO provider.",
				Required:    true,
			},
			"roles": schema.ListAttribute{
				Description: "List of Stacklet API roles automatically granted to SSO users in this group.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"account_group_uuids": schema.ListAttribute{
				Description: "List of account group UUIDs whose resources are visible to SSO users in this group.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *ssoGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ssoGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ssoGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		SSOGroupConfigs []ssoGroupConfig `graphql:"ssoGroupConfigs"`
	}

	err := d.client.Query(ctx, &query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSO groups, got error: %s", err))
		return
	}

	// Find our group
	var ourGroup *ssoGroupConfig
	for _, group := range query.SSOGroupConfigs {
		if group.Name == config.Name.ValueString() {
			ourGroup = &group
			break
		}
	}

	if ourGroup == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("SSO group with name %q not found", config.Name.ValueString()))
		return
	}

	// Generate a stable ID based on the name
	config.ID = types.StringValue(fmt.Sprintf("sso-group-%s", ourGroup.Name))
	config.Name = types.StringValue(ourGroup.Name)

	rolesValue, diags := types.ListValueFrom(ctx, types.StringType, ourGroup.Roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Roles = rolesValue

	accountGroupUUIDsValue, diags := types.ListValueFrom(ctx, types.StringType, ourGroup.AccountGroupUUIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.AccountGroupUUIDs = accountGroupUUIDsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
