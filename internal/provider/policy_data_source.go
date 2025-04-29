package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ datasource.DataSource = &policyDataSource{}
)

func NewPolicyDataSource() datasource.DataSource {
	return &policyDataSource{}
}

type policyDataSource struct {
	client *graphql.Client
}

type policyDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Version       types.Number `tfsdk:"version"`
}

func (d *policyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (d *policyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a policy by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the policy.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the policy.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the policy.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the policy (aws, azure, gcp, kubernetes, or tencentcloud).",
				Computed:    true,
			},
			"version": schema.NumberAttribute{
				Description: "The version of the policy.",
				Computed:    true,
			},
		},
	}
}

func (d *policyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *policyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data policyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		Policy struct {
			ID          string
			UUID        string
			Name        string
			Description string
			Provider    string
			Version     float64
		} `graphql:"policy(uuid: $uuid, name: $name)"`
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy, got error: %s", err))
		return
	}

	if query.Policy.UUID == "" {
		resp.Diagnostics.AddError("Not Found", "No policy found with the specified UUID or name")
		return
	}

	data.ID = types.StringValue(query.Policy.ID)
	data.UUID = types.StringValue(query.Policy.UUID)
	data.Name = types.StringValue(query.Policy.Name)
	data.Description = types.StringValue(query.Policy.Description)
	data.CloudProvider = types.StringValue(query.Policy.Provider)
	data.Version = types.NumberValue(big.NewFloat(query.Policy.Version))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
