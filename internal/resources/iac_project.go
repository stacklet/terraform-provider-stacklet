// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemadefault"
)

var (
	_ resource.Resource                = &iacProjectResource{}
	_ resource.ResourceWithConfigure   = &iacProjectResource{}
	_ resource.ResourceWithImportState = &iacProjectResource{}
)

type iacProjectResource struct {
	iacAPIResource
}

func (r *iacProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_project"
}

func (r *iacProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an IaC project.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the project.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID of the project credentials, if any.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"collections": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "The names of the policy collections used by the project.",
				Optional:    true,
				Computed:    true,
				Default:     schemadefault.EmptySetDefault(types.StringType),
			},
			"read_groups": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "The names of the groups with read access to the project.",
				Optional:    true,
				Computed:    true,
				Default:     schemadefault.EmptySetDefault(types.StringType),
			},
		},
	}
}

func (r *iacProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.IaCProjectResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := plan.ToAPIInput(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.api.Project.Create(ctx, plan.Name.ValueString(), input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, project)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iacProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.IaCProjectResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.api.Project.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(ctx, project)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *iacProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.IaCProjectResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := plan.ToAPIInput(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.api.Project.Update(ctx, plan.Name.ValueString(), input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, project)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iacProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.IaCProjectResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.Project.Delete(ctx, state.Name.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *iacProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"name"})
}
