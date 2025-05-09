package datasources

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &policyDataSource{}
)

func NewPolicyDataSource() datasource.DataSource {
	return &policyDataSource{}
}

type policyDataSource struct {
	api *api.API
}

func (d *policyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (d *policyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a policy, by UUID or name.",
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

	d.api = api.New(client)
}

func (d *policyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PolicyDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := d.api.Policy.Read(ctx, data.UUID.ValueString(), data.Name.ValueString())
	if err != nil {
		helpers.AddDiagError(resp.Diagnostics, err)
		return
	}

	if policy.UUID == "" {
		resp.Diagnostics.AddError("Not Found", "No policy found with the specified UUID or name")
		return
	}

	data.ID = types.StringValue(policy.ID)
	data.UUID = types.StringValue(policy.UUID)
	data.Name = types.StringValue(policy.Name)
	data.Description = tftypes.NullableString(policy.Description)
	data.CloudProvider = types.StringValue(policy.Provider)
	data.Version = types.NumberValue(big.NewFloat(policy.Version))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
