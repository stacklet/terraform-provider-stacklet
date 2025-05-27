package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
	"github.com/stacklet/terraform-provider-stacklet/schemavalidate"
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
				Description: "The UUID of the policy, alternative to the name.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy, alternative to the UUID.",
				Optional:    true,
			},
			"version": schema.Int32Attribute{
				Description: "The version policy. If not specified, the latest is used.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the policy.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the policy (aws, azure, gcp, kubernetes, or tencentcloud).",
				Computed:    true,
				Validators: []validator.String{
					schemavalidate.OneOfCloudProviders(),
				},
			},
			"category": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The list of categories the policy belongs to.",
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: "The policy mode.",
				Computed:    true,
			},
			"resource_type": schema.StringAttribute{
				Description: "The resource type that the policy applies to.",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "The path of the policy in the source repository.",
				Computed:    true,
			},
			"source_json": schema.StringAttribute{
				Description: "The policy source in JSON format.",
				Computed:    true,
			},
			"source_yaml": schema.StringAttribute{
				Description: "The policy source in YAML format.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system policy.",
				Computed:    true,
			},
			"unqualified_name": schema.StringAttribute{
				Description: "The policy name without namespace prefix.",
				Computed:    true,
			},
		},
	}
}

func (d *policyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *policyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PolicyDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.UUID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of UUID and name must be set")
	}

	policy, err := d.api.Policy.Read(ctx, data.UUID.ValueString(), data.Name.ValueString(), int(data.Version.ValueInt32()))
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(policy.ID)
	data.UUID = types.StringValue(policy.UUID)
	data.Name = types.StringValue(policy.Name)
	data.Description = tftypes.NullableString(policy.Description)
	data.CloudProvider = types.StringValue(policy.Provider)
	data.Version = types.Int32Value(int32(policy.Version))
	category, diag := types.ListValueFrom(ctx, types.StringType, policy.Category)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}
	data.Category = category
	data.Mode = types.StringValue(policy.Mode)
	data.ResourceType = types.StringValue(policy.ResourceType)
	data.Path = types.StringValue(policy.Path)
	data.SourceJSON = types.StringValue(policy.Source)
	data.SourceYAML = types.StringValue(policy.SourceYAML)
	data.System = types.BoolValue(policy.System)
	data.UnqualifiedName = types.StringValue(policy.UnqualifiedName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
