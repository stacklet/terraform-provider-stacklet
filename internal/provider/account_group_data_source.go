package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ datasource.DataSource = &accountGroupDataSource{}
)

func NewAccountGroupDataSource() datasource.DataSource {
	return &accountGroupDataSource{}
}

type accountGroupDataSource struct {
	client *graphql.Client
}

type accountGroupDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Regions       types.List   `tfsdk:"regions"`
}

func (d *accountGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group"
}

func (d *accountGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch an account group by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account group.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the account group.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the account group.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the account group.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account group.",
				Computed:    true,
			},
			"regions": schema.ListAttribute{
				Description: "The regions for the account group.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *accountGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		AccountGroup struct {
			ID          string
			UUID        string
			Name        string
			Description string
			Provider    string
			Regions     []string
		} `graphql:"accountGroup(uuid: $uuid, name: $name)"`
	}

	variables := map[string]interface{}{
		"uuid": (*string)(nil),
		"name": (*string)(nil),
	}

	if !data.UUID.IsNull() {
		variables["uuid"] = graphql.String(data.UUID.ValueString())
	}

	if !data.Name.IsNull() {
		variables["name"] = graphql.String(data.Name.ValueString())
	}

	err := d.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account group, got error: %s", err))
		return
	}

	if query.AccountGroup.UUID == "" {
		resp.Diagnostics.AddError("Not Found", "No account group found with the specified UUID or name")
		return
	}

	data.ID = types.StringValue(query.AccountGroup.ID)
	data.UUID = types.StringValue(query.AccountGroup.UUID)
	data.Name = types.StringValue(query.AccountGroup.Name)
	data.Description = types.StringValue(query.AccountGroup.Description)
	data.CloudProvider = types.StringValue(query.AccountGroup.Provider)
	regions := make([]attr.Value, len(query.AccountGroup.Regions))
	for i, region := range query.AccountGroup.Regions {
		regions[i] = types.StringValue(region)
	}
	data.Regions, _ = types.ListValue(types.StringType, regions)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
