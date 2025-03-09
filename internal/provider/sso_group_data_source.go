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
		Description: "Fetch an SSO group configuration by name.",
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
	var data ssoGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		SSOGroupConfigs []struct {
			Name              string
			Roles             []string
			AccountGroupUUIDs []string
		} `graphql:"ssoGroupConfigs"`
	}

	err := d.client.Query(ctx, &query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSO group configurations, got error: %s", err))
		return
	}

	// Find the matching group by name
	var matchingGroup *struct {
		Name              string
		Roles             []string
		AccountGroupUUIDs []string
	}
	for _, group := range query.SSOGroupConfigs {
		if group.Name == data.Name.ValueString() {
			matchingGroup = &group
			break
		}
	}

	if matchingGroup == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No SSO group found with name: %s", data.Name.ValueString()))
		return
	}

	// Generate a stable ID based on the name
	data.ID = types.StringValue(fmt.Sprintf("sso-group-%s", matchingGroup.Name))
	data.Name = types.StringValue(matchingGroup.Name)

	roles, diags := types.ListValueFrom(ctx, types.StringType, matchingGroup.Roles)
	resp.Diagnostics.Append(diags...)
	data.Roles = roles

	accountGroupUUIDs, diags := types.ListValueFrom(ctx, types.StringType, matchingGroup.AccountGroupUUIDs)
	resp.Diagnostics.Append(diags...)
	data.AccountGroupUUIDs = accountGroupUUIDs

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
