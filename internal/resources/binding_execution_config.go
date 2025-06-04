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

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &bindingExecutionConfigResource{}
	_ resource.ResourceWithConfigure   = &bindingExecutionConfigResource{}
	_ resource.ResourceWithImportState = &bindingExecutionConfigResource{}
)

func NewBindingExecutionConfigResource() resource.Resource {
	return &bindingExecutionConfigResource{}
}

type bindingExecutionConfigResource struct {
	api *api.API
}

func (r *bindingExecutionConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding_execution_config"
}

func (r *bindingExecutionConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage execution configuration for a binding.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL ID of the binding the configuration is for.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"binding_uuid": schema.StringAttribute{
				Description: "The UUID of the binding the configuration is for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
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
	}
}

func (r *bindingExecutionConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *bindingExecutionConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.BindingExecutionConfigResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.BindingExecutionConfigUpdateInput{
		BindingUUID: plan.BindingUUID.ValueString(),
		ExecutionConfig: api.BindingExecutionConfig{
			DryRun: &api.BindingExecutionConfigDryRun{
				Default: plan.DryRun.ValueBool(),
			},
			Variables: plan.Variables.ValueStringPointer(),
		},
	}
	executionConfig, err := r.api.BindingExecutionConfig.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateBindingExecutionConfigModel(&plan, executionConfig))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bindingExecutionConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.BindingExecutionConfigResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	executionConfig, err := r.api.BindingExecutionConfig.Read(ctx, state.BindingUUID.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	r.updateBindingExecutionConfigModel(&state, executionConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bindingExecutionConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.BindingExecutionConfigResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.BindingExecutionConfigUpdateInput{
		BindingUUID: plan.BindingUUID.ValueString(),
		ExecutionConfig: api.BindingExecutionConfig{
			DryRun: &api.BindingExecutionConfigDryRun{
				Default: plan.DryRun.ValueBool(),
			},
			Variables: plan.Variables.ValueStringPointer(),
		},
	}
	executionConfig, err := r.api.BindingExecutionConfig.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateBindingExecutionConfigModel(&plan, executionConfig))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bindingExecutionConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.BindingExecutionConfigResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete means setting the execution config to the empty value
	input := api.BindingExecutionConfigUpdateInput{
		BindingUUID:     state.BindingUUID.ValueString(),
		ExecutionConfig: api.BindingExecutionConfig{},
	}
	if _, err := r.api.BindingExecutionConfig.Upsert(ctx, input); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	}
}

func (r *bindingExecutionConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("binding_uuid"), req.ID)...)
}

func (r bindingExecutionConfigResource) updateBindingExecutionConfigModel(m *models.BindingExecutionConfigResource, executionConfig *api.BindingExecutionConfig) diag.Diagnostic {
	m.ID = m.BindingUUID // since it's not a separate GraphQL entity, use the binding UUID as ID too
	m.DryRun = types.BoolValue(executionConfig.DryRunDefault())
	variablesString, err := tftypes.JSONString(executionConfig.Variables)
	if err != nil {
		return diag.NewErrorDiagnostic("Invalid content for variables", err.Error())
	}

	// the API returns an empty dict for null or empty string, don't modify the
	// expected value in that case
	if variablesString.ValueString() != "{}" {
		m.Variables = variablesString
	}
	return nil
}
