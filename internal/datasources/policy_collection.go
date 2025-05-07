package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ datasource.DataSource = &policyCollectionDataSource{}
)

func NewPolicyCollectionDataSource() datasource.DataSource {
	return &policyCollectionDataSource{}
}

type policyCollectionDataSource struct {
	client *graphql.Client
}

type policyCollectionDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	AutoUpdate    types.Bool   `tfsdk:"auto_update"`
	System        types.Bool   `tfsdk:"system"`
	Repository    types.String `tfsdk:"repository"`
}

func (d *policyCollectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection"
}

func (d *policyCollectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a policy collection by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the policy collection.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy collection.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the policy collection.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the policy collection.",
				Computed:    true,
			},
			"auto_update": schema.BoolAttribute{
				Description: "Whether the policy collection automatically updates policy versions.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system policy collection.",
				Computed:    true,
			},
			"repository": schema.StringAttribute{
				Description: "The repository URL if this collection was created from a repo control file.",
				Computed:    true,
			},
		},
	}
}

func (d *policyCollectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *policyCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data policyCollectionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		PolicyCollection struct {
			ID          string
			UUID        string
			Name        string
			Description string
			Provider    string
			AutoUpdate  bool
			System      bool
			Repository  string
		} `graphql:"policyCollection(uuid: $uuid, name: $name)"`
	}

	variables := map[string]any{
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy collection, got error: %s", err))
		return
	}

	if query.PolicyCollection.UUID == "" {
		resp.Diagnostics.AddError("Not Found", "No policy collection found with the specified UUID or name")
		return
	}

	data.ID = types.StringValue(query.PolicyCollection.ID)
	data.UUID = types.StringValue(query.PolicyCollection.UUID)
	data.Name = types.StringValue(query.PolicyCollection.Name)
	data.Description = types.StringValue(query.PolicyCollection.Description)
	data.CloudProvider = types.StringValue(query.PolicyCollection.Provider)
	data.AutoUpdate = types.BoolValue(query.PolicyCollection.AutoUpdate)
	data.System = types.BoolValue(query.PolicyCollection.System)
	data.Repository = types.StringValue(query.PolicyCollection.Repository)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
