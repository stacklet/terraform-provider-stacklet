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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemavalidate"
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
			"dry_run": schema.BoolAttribute{
				Description: "Whether the binding is run in with action disabled (in information mode).",
				Optional:    true,
			},
			"resource_limits": schema.SingleNestedAttribute{
				Description: "Default resource limits for binding execution.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"max_count": schema.Int32Attribute{
						Description: "Max count of affected resources.",
						Optional:    true,
					},
					"max_percentage": schema.Float32Attribute{
						Description: "Max percentage of affected resources.",
						Optional:    true,
					},
					"requires_both": schema.BoolAttribute{
						Description: "If set, only applies limits when both thresholds are exceeded.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
				Validators: []validator.Object{
					bindingResourceLimitsValidator{},
				},
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
		Blocks: map[string]schema.Block{
			"policy_resource_limit": schema.ListNestedBlock{
				Description: "Per-policy overrides for resource limits for binding execution. Map keys are policy unqualified names.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"policy_name": schema.StringAttribute{
							Description: "Unqualified name of the policy for the limit override.",
							Required:    true,
						},
						"max_count": schema.Int32Attribute{
							Description: "Max count of affected resources.",
							Optional:    true,
						},
						"max_percentage": schema.Float32Attribute{
							Description: "Max percentage of affected resources.",
							Optional:    true,
						},
						"requires_both": schema.BoolAttribute{
							Description: "If set, only applies limits when both thresholds are exceeded.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
					Validators: []validator.Object{
						bindingResourcePolicyLimitsValidator{},
					},
				},
				Validators: []validator.List{
					schemavalidate.UniqueStringAttribute("policy_name"),
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

	executionConfig, diags := r.getExecutionConfig(ctx, plan, config.SecurityContextWO.ValueStringPointer())
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

	resp.Diagnostics.Append(plan.Update(ctx, binding)...)
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

	resp.Diagnostics.Append(state.Update(ctx, binding)...)
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

	var securityContextString *string
	if state.SecurityContextWOVersion == plan.SecurityContextWOVersion {
		// if no change happened, send the value we got from the API as a
		// result of the previous change. Not sending a value makes the API
		// unset it.
		securityContextString = state.SecurityContext.ValueStringPointer()
	} else {
		securityContextString = config.SecurityContextWO.ValueStringPointer()
	}

	executionConfig, diags := r.getExecutionConfig(ctx, plan, securityContextString)
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

	resp.Diagnostics.Append(plan.Update(ctx, binding)...)
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

func (r bindingResource) getExecutionConfig(ctx context.Context, plan models.BindingResource, securityContextString *string) (api.BindingExecutionConfig, diag.Diagnostics) {
	var dryRun *api.BindingExecutionConfigDryRun
	if !plan.DryRun.IsNull() {
		dryRun = &api.BindingExecutionConfigDryRun{Default: plan.DryRun.ValueBool()}
	}

	var securityContext *api.BindingExecutionConfigSecurityContext
	if securityContextString != nil {
		securityContext = &api.BindingExecutionConfigSecurityContext{Default: *securityContextString}
	}

	var defaultResourceLimits *api.BindingExecutionConfigResourceLimit
	if !plan.ResourceLimits.IsNull() {
		var defLimitsObj models.BindingExecutionConfigResourceLimit
		if diags := plan.ResourceLimits.As(ctx, &defLimitsObj, ObjectAsOptions); diags.HasError() {
			return api.BindingExecutionConfig{}, diags
		}
		defaultResourceLimits = &api.BindingExecutionConfigResourceLimit{
			MaxCount:      defLimitsObj.MaxCount.ValueInt32Pointer(),
			MaxPercentage: defLimitsObj.MaxPercentage.ValueFloat32Pointer(),
			RequiresBoth:  defLimitsObj.RequiresBoth.ValueBool(),
		}
	}

	policyResourceLimits := make([]api.BindingExecutionConfigResourceLimitsPolicyOverrides, 0)
	if !plan.PolicyResourceLimits.IsNull() {
		for i, elem := range plan.PolicyResourceLimits.Elements() {
			resourceLimit, ok := elem.(types.Object)
			if !ok {
				var diags diag.Diagnostics
				diags.AddAttributeError(
					path.Root("policy_resource_limit").AtListIndex(i),
					"Invalid limits",
					"Limits block is invalid,",
				)
				return api.BindingExecutionConfig{}, diags
			}
			var limitsObj models.BindingExecutionConfigPolicyResourceLimit
			if diags := resourceLimit.As(ctx, &limitsObj, ObjectAsOptions); diags.HasError() {
				return api.BindingExecutionConfig{}, diags
			}

			policyResourceLimits = append(
				policyResourceLimits,
				api.BindingExecutionConfigResourceLimitsPolicyOverrides{
					PolicyName: limitsObj.PolicyName.ValueString(),
					Limit: api.BindingExecutionConfigResourceLimit{
						MaxCount:      limitsObj.MaxCount.ValueInt32Pointer(),
						MaxPercentage: limitsObj.MaxPercentage.ValueFloat32Pointer(),
						RequiresBoth:  limitsObj.RequiresBoth.ValueBool(),
					},
				},
			)
		}
	}

	return api.BindingExecutionConfig{
		DryRun: dryRun,
		ResourceLimits: &api.BindingExecutionConfigResourceLimits{
			Default:         defaultResourceLimits,
			PolicyOverrides: policyResourceLimits,
		},
		SecurityContext: securityContext,
		Variables:       plan.Variables.ValueStringPointer(),
	}, nil
}

func bindingExecutionConfigLimitValidateObject(obj models.BindingExecutionConfigResourceLimit, objPath path.Path, diags *diag.Diagnostics) {
	if obj.MaxCount.IsNull() && obj.MaxPercentage.IsNull() {
		diags.AddAttributeError(
			objPath,
			"Invalid value",
			"At least one of `max_path` and `max_percentage` must be set",
		)
		return
	}

	if obj.RequiresBoth.IsNull() {
		return
	}
	if obj.RequiresBoth.ValueBool() && (obj.MaxCount.IsNull() || obj.MaxPercentage.IsNull()) {
		diags.AddAttributeError(
			objPath.AtName("requires_both"),
			"Invalid value",
			"The attribute can be set to true only if both limits are set",
		)
	}
}

type bindingResourceLimitsValidator struct{}

func (m bindingResourceLimitsValidator) Description(ctx context.Context) string {
	return "Check that resource limits for bindings are properly configured."
}

func (m bindingResourceLimitsValidator) MarkdownDescription(ctx context.Context) string {
	return "Check that resource limits for bindings are properly configured."
}

func (m bindingResourceLimitsValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var obj models.BindingExecutionConfigResourceLimit
	if diags := req.ConfigValue.As(ctx, &obj, ObjectAsOptions); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	bindingExecutionConfigLimitValidateObject(obj, req.Path, &resp.Diagnostics)
}

type bindingResourcePolicyLimitsValidator struct{}

func (m bindingResourcePolicyLimitsValidator) Description(ctx context.Context) string {
	return "Check that policy resource limits overrides for bindings are properly configured."
}

func (m bindingResourcePolicyLimitsValidator) MarkdownDescription(ctx context.Context) string {
	return "Check that policy resource limits overrides for bindings are properly configured."
}

func (m bindingResourcePolicyLimitsValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var obj models.BindingExecutionConfigPolicyResourceLimit
	if diags := req.ConfigValue.As(ctx, &obj, ObjectAsOptions); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	bindingExecutionConfigLimitValidateObject(obj.BindingExecutionConfigResourceLimit, req.Path, &resp.Diagnostics)
}
