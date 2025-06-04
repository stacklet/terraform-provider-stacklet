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
					"security_context": schema.StringAttribute{
						Description: "The binding execution security context.",
						Optional:    true,
						Computed:    true,
					},
					"security_context_wo": schema.StringAttribute{
						Description: "The input value for the security context for the execution configuration.",
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
					},
					"security_context_wo_version": schema.StringAttribute{
						Description: "The version for the security context. Must be changed to update security_context_wo.",
						Optional:    true,
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

	executionConfig, diags := r.getCreateExecutionConfig(ctx, plan.ExecutionConfig, config.ExecutionConfig)
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
	var plan, state, config models.BindingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	executionConfig, diags := r.getUpdateExecutionConfig(ctx, plan.ExecutionConfig, state.ExecutionConfig, config.ExecutionConfig)
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

func (r bindingResource) getCreateExecutionConfig(ctx context.Context, plan, config types.Object) (api.BindingExecutionConfig, diag.Diagnostics) {
	var executionConfig api.BindingExecutionConfig
	var diags diag.Diagnostics

	if plan.IsNull() {
		return executionConfig, diags
	}

	var planObj, configObj models.BindingResourceExecutionConfig
	diags.Append(plan.As(ctx, &planObj, basetypes.ObjectAsOptions{})...)
	diags.Append(config.As(ctx, &configObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return executionConfig, diags
	}

	var dryRun *api.BindingExecutionConfigDryRun
	if !planObj.DryRun.IsNull() {
		dryRun = &api.BindingExecutionConfigDryRun{Default: planObj.DryRun.ValueBool()}
	}
	return api.BindingExecutionConfig{
		DryRun:          dryRun,
		SecurityContext: &api.BindingExecutionConfigSecurityContext{Default: configObj.SecurityContextWO.ValueString()},
		Variables:       planObj.Variables.ValueStringPointer(),
	}, diags
}

func (r bindingResource) getUpdateExecutionConfig(ctx context.Context, plan, state, config types.Object) (api.BindingExecutionConfig, diag.Diagnostics) {
	var executionConfig api.BindingExecutionConfig
	var diags diag.Diagnostics

	if plan.IsNull() {
		return executionConfig, diags
	}

	m, d := r.getExecutionConfigModels(ctx, plan, state, config)
	diags.Append(d...)
	if diags.HasError() {
		return executionConfig, diags
	}

	var dryRun *api.BindingExecutionConfigDryRun
	if !m.Plan.DryRun.IsNull() {
		dryRun = &api.BindingExecutionConfigDryRun{Default: m.Plan.DryRun.ValueBool()}
	}
	var securityContext string
	if m.State.SecurityContextWOVersion == m.Plan.SecurityContextWOVersion {
		// if no change happened, send the value we got from the API as a
		// result of the previous change. Not sending a value makes the API
		// unset it.
		securityContext = m.State.SecurityContext.ValueString()
	} else {
		securityContext = m.Config.SecurityContextWO.ValueString()
	}

	return api.BindingExecutionConfig{
		DryRun:          dryRun,
		SecurityContext: &api.BindingExecutionConfigSecurityContext{Default: securityContext},
		Variables:       m.Plan.Variables.ValueStringPointer(),
	}, diags
}

type bindingExecutionConfigModels struct {
	Plan   models.BindingResourceExecutionConfig
	State  models.BindingResourceExecutionConfig
	Config models.BindingResourceExecutionConfig
}

func (r bindingResource) getExecutionConfigModels(ctx context.Context, plan, state, config types.Object) (bindingExecutionConfigModels, diag.Diagnostics) {
	var models bindingExecutionConfigModels
	var diags diag.Diagnostics
	if !plan.IsNull() {
		diags.Append(plan.As(ctx, &models.Plan, basetypes.ObjectAsOptions{})...)
	}
	if !state.IsNull() {
		diags.Append(state.As(ctx, &models.State, basetypes.ObjectAsOptions{})...)
	}
	if !config.IsNull() {
		diags.Append(config.As(ctx, &models.Config, basetypes.ObjectAsOptions{})...)
	}
	return models, diags
}

func (r bindingResource) updateBindingModel(ctx context.Context, m, config *models.BindingResource, binding *api.Binding) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(binding.ID)
	m.UUID = types.StringValue(binding.UUID)
	m.Name = types.StringValue(binding.Name)
	m.Description = tftypes.NullableString(binding.Description)
	m.AutoDeploy = types.BoolValue(binding.AutoDeploy)
	m.Schedule = tftypes.NullableString(binding.Schedule)
	m.AccountGroupUUID = types.StringValue(binding.AccountGroup.UUID)
	m.PolicyCollectionUUID = types.StringValue(binding.PolicyCollection.UUID)
	m.System = types.BoolValue(binding.System)

	executionConfig, d := tftypes.ObjectValue(
		ctx,
		&(binding.ExecutionConfig),
		func() (*models.BindingResourceExecutionConfig, error) {
			empty := api.BindingExecutionConfig{}
			if binding.ExecutionConfig == empty {
				// requested config was null
				if config == nil || config.ExecutionConfig.IsNull() {
					return nil, nil
				}
				// requested config was empty
				return &models.BindingResourceExecutionConfig{
					DryRun:          types.BoolValue(false),
					SecurityContext: types.StringNull(),
					Variables:       types.StringNull(),
				}, nil
			}

			var curModel models.BindingResourceExecutionConfig
			if !m.ExecutionConfig.IsNull() {
				diags.Append(m.ExecutionConfig.As(ctx, &curModel, basetypes.ObjectAsOptions{})...)
			}

			variablesString, err := tftypes.JSONString(binding.ExecutionConfig.Variables)
			if err != nil {
				return nil, err
			}
			// the API returns an empty dict for null or empty string, don't
			// modify the expected value in that case
			if variablesString.ValueString() == "{}" {
				variablesString = curModel.Variables
			}

			return &models.BindingResourceExecutionConfig{
				DryRun:                   types.BoolValue(binding.ExecutionConfig.DryRunDefault()),
				SecurityContext:          tftypes.NullableString(binding.ExecutionConfig.SecurityContextDefault()),
				SecurityContextWOVersion: curModel.SecurityContextWOVersion,
				Variables:                variablesString,
			}, nil
		},
	)
	m.ExecutionConfig = executionConfig
	diags.Append(d...)
	return diags
}
