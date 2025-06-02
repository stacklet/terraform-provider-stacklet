// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &bindingResource{}
	_ resource.ResourceWithConfigure   = &bindingResource{}
	_ resource.ResourceWithImportState = &bindingResource{}
)

func NewBindingResource() resource.Resource {
	return &bindingResource{}
}

type bindingResource struct {
	api *api.API
}

func (r *bindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding"
}

func (r *bindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a binding between an account group and a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the binding.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the binding.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the binding.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the binding.",
				Optional:    true,
			},
			"auto_deploy": schema.BoolAttribute{
				Description: "Whether the binding should automatically deploy when the policy collection changes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"schedule": schema.StringAttribute{
				Description: "The schedule for the binding (e.g., 'rate(1 hour)', 'rate(2 hours)', or cron expression).",
				Optional:    true,
			},
			"account_group_uuid": schema.StringAttribute{
				Description: "The UUID of the account group this binding applies to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_collection_uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection this binding applies.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system": schema.BoolAttribute{
				Description: "Whether the binding is a system one. Always false for resources.",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"execution_config": schema.SingleNestedAttribute{
				Description: "Binding execution configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"dry_run": schema.BoolAttribute{
						Description: "Whether the binding is run in with action disabled (in information mode).",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"variables": schema.StringAttribute{
						Description: "JSON-encoded dictionary of values used for policy templating.",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r *bindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *bindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.BindingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	executionConfig, diags := r.getExecutionConfig(ctx, plan.ExecutionConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.BindingCreateInput{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueStringPointer(),
		AutoDeploy:           plan.AutoDeploy.ValueBool(),
		Schedule:             plan.Schedule.ValueStringPointer(),
		ExecutionConfig:      executionConfig,
		AccountGroupUUID:     plan.AccountGroupUUID.ValueString(),
		PolicyCollectionUUID: plan.PolicyCollectionUUID.ValueString(),
		Deploy:               true,
	}

	binding, err := r.api.Binding.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateBindingModel(ctx, &plan, &config, binding)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.BindingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	binding, err := r.api.Binding.Read(ctx, state.UUID.ValueString(), "")
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateBindingModel(ctx, &state, nil, binding)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, config models.BindingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	executionConfig, diags := r.getExecutionConfig(ctx, plan.ExecutionConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	input := api.BindingUpdateInput{
		UUID:            plan.UUID.ValueString(),
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueStringPointer(),
		AutoDeploy:      plan.AutoDeploy.ValueBool(),
		Schedule:        plan.Schedule.ValueStringPointer(),
		ExecutionConfig: executionConfig,
	}

	binding, err := r.api.Binding.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateBindingModel(ctx, &plan, &config, binding)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.BindingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.Binding.Delete(ctx, state.UUID.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *bindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), req.ID)...)
}

func (r bindingResource) getExecutionConfig(ctx context.Context, planExecutionConfig types.Object) (api.BindingExecutionConfig, diag.Diagnostics) {
	var executionConfig api.BindingExecutionConfig
	var diags diag.Diagnostics

	var configObj models.BindingExecutionConfig
	if !planExecutionConfig.IsNull() {
		diags.Append(planExecutionConfig.As(ctx, &configObj, basetypes.ObjectAsOptions{})...)

		var dryRun *api.BindingExecutionConfigDryRun
		if !configObj.DryRun.IsNull() {
			dryRun = &api.BindingExecutionConfigDryRun{Default: configObj.DryRun.ValueBool()}
		}
		executionConfig = api.BindingExecutionConfig{
			DryRun:    dryRun,
			Variables: configObj.Variables.ValueStringPointer(),
		}
	}

	return executionConfig, diags
}

func (r bindingResource) updateBindingModel(ctx context.Context, m *models.BindingResource, config *models.BindingResource, binding *api.Binding) diag.Diagnostics {
	m.ID = types.StringValue(binding.ID)
	m.UUID = types.StringValue(binding.UUID)
	m.Name = types.StringValue(binding.Name)
	m.Description = tftypes.NullableString(binding.Description)
	m.AutoDeploy = types.BoolValue(binding.AutoDeploy)
	m.Schedule = tftypes.NullableString(binding.Schedule)
	m.AccountGroupUUID = types.StringValue(binding.AccountGroup.UUID)
	m.PolicyCollectionUUID = types.StringValue(binding.PolicyCollection.UUID)
	m.System = types.BoolValue(binding.System)

	executionConfig, diags := tftypes.ObjectValue(
		ctx,
		&(binding.ExecutionConfig),
		func() (*models.BindingExecutionConfig, error) {
			empty := api.BindingExecutionConfig{}
			if binding.ExecutionConfig == empty {
				// requested config was null
				if config == nil || config.ExecutionConfig.IsNull() {
					return nil, nil
				}
				// requested config was empty
				return &models.BindingExecutionConfig{
					DryRun:    types.BoolValue(false),
					Variables: types.StringNull(),
				}, nil
			}
			variables, err := tftypes.JSONString(binding.ExecutionConfig.Variables)
			if err != nil {
				return nil, err
			}
			dryRun := types.BoolValue(false)
			if binding.ExecutionConfig.DryRun != nil {
				dryRun = types.BoolValue(binding.ExecutionConfig.DryRun.Default)
			}

			return &models.BindingExecutionConfig{
				DryRun:    dryRun,
				Variables: variables,
			}, nil
		},
	)
	m.ExecutionConfig = executionConfig
	return diags
}
