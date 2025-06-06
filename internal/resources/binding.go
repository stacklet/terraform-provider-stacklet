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
			"dry_run": schema.BoolAttribute{
				Description: "Whether the binding is run in with action disabled (in information mode).",
				Optional:    true,
			},
			"default_resource_limits": schema.SingleNestedAttribute{
				Description: "Default limits for binding execution.",
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

	executionConfig, diags := r.getCreateExecutionConfig(ctx, plan, config)
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

	resp.Diagnostics.Append(r.updateBindingModel(ctx, &plan, binding)...)
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

	resp.Diagnostics.Append(r.updateBindingModel(ctx, &state, binding)...)
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

	executionConfig, diags := r.getUpdateExecutionConfig(ctx, plan, state, config)
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

	resp.Diagnostics.Append(r.updateBindingModel(ctx, &plan, binding)...)
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

func (r bindingResource) updateBindingModel(ctx context.Context, m *models.BindingResource, binding *api.Binding) diag.Diagnostics {
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
	m.DryRun = tftypes.NullableBool(binding.DryRun())
	m.SecurityContext = tftypes.NullableString(binding.SecurityContext())

	variablesString, err := tftypes.JSONString(binding.ExecutionConfig.Variables)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("Invalid content for variables", err.Error()))
		return diags
	}
	// the API returns an empty dict for null or empty string, don't
	// modify the expected value in that case
	if variablesString.ValueString() == "{}" {
		variablesString = m.Variables
	}
	m.Variables = variablesString

	var defaultLimits types.Object
	if binding.DefaultResourceLimits() == nil && !m.DefaultResourceLimits.IsNull() {
		var d diag.Diagnostics
		// if default resource limits are set in the config but empty, apply default
		def := models.BindingExecutionConfigResourceLimit{RequiresBoth: types.BoolValue(false)}
		defaultLimits, d = basetypes.NewObjectValueFrom(ctx, def.AttributeTypes(), &def)
		diags.Append(d...)
	} else {
		var d diag.Diagnostics
		defaultLimits, d = tftypes.ObjectValue(
			ctx,
			binding.DefaultResourceLimits(),
			func() (*models.BindingExecutionConfigResourceLimit, diag.Diagnostics) {
				l := binding.DefaultResourceLimits()
				return &models.BindingExecutionConfigResourceLimit{
					MaxCount:      tftypes.NullableInt(l.MaxCount),
					MaxPercentage: tftypes.NullableFloat(l.MaxPercentage),
					RequiresBoth:  types.BoolValue(l.RequiresBoth),
				}, nil
			},
		)
		diags.Append(d...)
	}
	m.DefaultResourceLimits = defaultLimits
	return diags
}

func (r bindingResource) getCreateExecutionConfig(ctx context.Context, plan, config models.BindingResource) (api.BindingExecutionConfig, diag.Diagnostics) {
	var dryRun *api.BindingExecutionConfigDryRun
	if !plan.DryRun.IsNull() {
		dryRun = &api.BindingExecutionConfigDryRun{Default: plan.DryRun.ValueBool()}
	}
	var securityContext *api.BindingExecutionConfigSecurityContext
	if !config.SecurityContextWO.IsNull() {
		securityContext = &api.BindingExecutionConfigSecurityContext{Default: config.SecurityContextWO.ValueString()}
	}
	var defaultResourceLimits *api.BindingExecutionConfigResourceLimit
	if !plan.DefaultResourceLimits.IsNull() {
		var defLimitsObj models.BindingExecutionConfigResourceLimit
		if diags := plan.DefaultResourceLimits.As(ctx, &defLimitsObj, basetypes.ObjectAsOptions{}); diags.HasError() {
			return api.BindingExecutionConfig{}, diags
		}
		defaultResourceLimits = &api.BindingExecutionConfigResourceLimit{
			MaxCount:      api.NullableInt(defLimitsObj.MaxCount),
			MaxPercentage: defLimitsObj.MaxPercentage.ValueFloat32Pointer(),
			RequiresBoth:  defLimitsObj.RequiresBoth.ValueBool(),
		}
	}

	return api.BindingExecutionConfig{
		DryRun: dryRun,
		ResourceLimits: &api.BindingExecutionConfigResourceLimits{
			Default: defaultResourceLimits,
		},
		SecurityContext: securityContext,
		Variables:       plan.Variables.ValueStringPointer(),
	}, nil
}

func (r bindingResource) getUpdateExecutionConfig(ctx context.Context, plan, state, config models.BindingResource) (api.BindingExecutionConfig, diag.Diagnostics) {
	var dryRun *api.BindingExecutionConfigDryRun
	if !plan.DryRun.IsNull() {
		dryRun = &api.BindingExecutionConfigDryRun{Default: plan.DryRun.ValueBool()}
	}
	var securityContextString *string
	var securityContext *api.BindingExecutionConfigSecurityContext
	if state.SecurityContextWOVersion == plan.SecurityContextWOVersion {
		// if no change happened, send the value we got from the API as a
		// result of the previous change. Not sending a value makes the API
		// unset it.
		securityContextString = state.SecurityContext.ValueStringPointer()
	} else {
		securityContextString = config.SecurityContextWO.ValueStringPointer()
	}
	if securityContextString != nil {
		securityContext = &api.BindingExecutionConfigSecurityContext{Default: *securityContextString}
	}

	var defaultResourceLimits *api.BindingExecutionConfigResourceLimit
	if !plan.DefaultResourceLimits.IsNull() {
		var defLimitsObj models.BindingExecutionConfigResourceLimit
		if diags := plan.DefaultResourceLimits.As(ctx, &defLimitsObj, basetypes.ObjectAsOptions{}); diags.HasError() {
			return api.BindingExecutionConfig{}, diags
		}
		defaultResourceLimits = &api.BindingExecutionConfigResourceLimit{
			MaxCount:      api.NullableInt(defLimitsObj.MaxCount),
			MaxPercentage: defLimitsObj.MaxPercentage.ValueFloat32Pointer(),
			RequiresBoth:  defLimitsObj.RequiresBoth.ValueBool(),
		}
	}

	return api.BindingExecutionConfig{
		DryRun: dryRun,
		ResourceLimits: &api.BindingExecutionConfigResourceLimits{
			Default: defaultResourceLimits,
		},
		SecurityContext: securityContext,
		Variables:       plan.Variables.ValueStringPointer(),
	}, nil
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
	if diags := req.ConfigValue.As(ctx, &obj, basetypes.ObjectAsOptions{}); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if obj.RequiresBoth.IsNull() {
		return
	}
	if obj.RequiresBoth.ValueBool() && (obj.MaxCount.IsNull() || obj.MaxPercentage.IsNull()) {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("required_both"),
			"Invalid value",
			"The attribute can be set to true only if both limits are set",
		)
	}
}
