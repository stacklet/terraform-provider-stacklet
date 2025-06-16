// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &bindingDataSource{}
)

func NewBindingDataSource() datasource.DataSource {
	return &bindingDataSource{}
}

type bindingDataSource struct {
	api *api.API
}

func (d *bindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding"
}

func (d *bindingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a binding by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the binding.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the binding.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the binding.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the binding.",
				Computed:    true,
			},
			"auto_deploy": schema.BoolAttribute{
				Description: "Whether the binding automatically deploys when the policy collection changes.",
				Computed:    true,
			},
			"schedule": schema.StringAttribute{
				Description: "The schedule for the binding (e.g., 'rate(1 hour)', 'rate(2 hours)', or cron expression).",
				Computed:    true,
			},
			"account_group_uuid": schema.StringAttribute{
				Description: "The UUID of the account group this binding applies to.",
				Computed:    true,
			},
			"policy_collection_uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection this binding applies.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system binding.",
				Computed:    true,
			},
			"dry_run": schema.BoolAttribute{
				Description: "Whether the binding is run in with action disabled (in information mode).",
				Optional:    true,
				Computed:    true,
			},
			"resource_limits": schema.SingleNestedAttribute{
				Description: "Default resource limits for binding execution.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"max_count": schema.Int32Attribute{
						Description: "Max count of affected resources.",
						Optional:    true,
						Computed:    true,
					},
					"max_percentage": schema.Float32Attribute{
						Description: "Max percentage of affected resources.",
						Optional:    true,
						Computed:    true,
					},
					"requires_both": schema.BoolAttribute{
						Description: "If set, only applies limits when both thresholds are exceeded.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"security_context": schema.StringAttribute{
				Description: "The binding execution security context.",
				Optional:    true,
				Computed:    true,
			},
			"variables": schema.StringAttribute{
				Description: "JSON-encoded dictionary of values used for policy templating.",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"policy_resource_limit": schema.ListNestedBlock{
				Description: "Per-policy overrides for resource limits for binding execution. Map keys are policy unqualified names.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"policy_name": schema.StringAttribute{
							Description: "Unqualified name of the policy for the limit override.",
							Computed:    true,
						},
						"max_count": schema.Int32Attribute{
							Description: "Max count of affected resources.",
							Optional:    true,
							Computed:    true,
						},
						"max_percentage": schema.Float32Attribute{
							Description: "Max percentage of affected resources.",
							Optional:    true,
							Computed:    true,
						},
						"requires_both": schema.BoolAttribute{
							Description: "If set, only applies limits when both thresholds are exceeded.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *bindingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *bindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.BindingDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	binding, err := d.api.Binding.Read(ctx, data.UUID.ValueString(), data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(binding.ID)
	data.UUID = types.StringValue(binding.UUID)
	data.Name = types.StringValue(binding.Name)
	data.Description = tftypes.NullableString(binding.Description)
	data.AutoDeploy = types.BoolValue(binding.AutoDeploy)
	data.Schedule = tftypes.NullableString(binding.Schedule)
	data.AccountGroupUUID = types.StringValue(binding.AccountGroup.UUID)
	data.PolicyCollectionUUID = types.StringValue(binding.PolicyCollection.UUID)
	data.System = types.BoolValue(binding.System)
	data.DryRun = tftypes.NullableBool(binding.DryRun())
	data.SecurityContext = tftypes.NullableString(binding.SecurityContext())

	variablesString, err := tftypes.JSONString(binding.ExecutionConfig.Variables)
	if err != nil {
		errors.AddDiagAttributeError(&resp.Diagnostics, "variables", err)
		return
	}
	data.Variables = variablesString

	defLimit := binding.DefaultResourceLimits()
	defaultLimits, diags := tftypes.ObjectValue(
		ctx,
		defLimit,
		func() (*models.BindingExecutionConfigResourceLimit, diag.Diagnostics) {
			return &models.BindingExecutionConfigResourceLimit{
				MaxCount:      tftypes.NullableInt(defLimit.MaxCount),
				MaxPercentage: tftypes.NullableFloat(defLimit.MaxPercentage),
				RequiresBoth:  types.BoolValue(defLimit.RequiresBoth),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ResourceLimits = defaultLimits

	policyLimits, diags := tftypes.ObjectList[models.BindingExecutionConfigPolicyResourceLimit](
		binding.PolicyResourceLimits(),
		func(entry api.BindingExecutionConfigResourceLimitsPolicyOverrides) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"policy_name":    types.StringValue(entry.PolicyName),
				"max_count":      tftypes.NullableInt(entry.Limit.MaxCount),
				"max_percentage": tftypes.NullableFloat(entry.Limit.MaxPercentage),
				"requires_both":  types.BoolValue(entry.Limit.RequiresBoth),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.PolicyResourceLimits = policyLimits

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
