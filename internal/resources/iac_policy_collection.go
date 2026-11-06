// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemadefault"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemavalidate"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &iacPolicyCollectionResource{}
	_ resource.ResourceWithConfigure   = &iacPolicyCollectionResource{}
	_ resource.ResourceWithImportState = &iacPolicyCollectionResource{}
)

// iacPolicyCollectionResource defines the resource implementation.
type iacPolicyCollectionResource struct {
	iacAPIResource
}

func (r *iacPolicyCollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_policy_collection"
}

func (r *iacPolicyCollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Stacklet IaC policy collection, a set of IaC policies to evaluate. Deleting a policy collection also removes it from all projects.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique name of the policy collection. It can't end with \"-mirror\", as that's reserved for collections mirroring policy sources.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					schemavalidate.NoSuffix("-mirror"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policies": schema.SetAttribute{
				Description: "The names of policies in the collection.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     schemadefault.EmptySetDefault(types.StringType),
			},
			"is_default": schema.BoolAttribute{
				Description: "Whether the collection is a default, applied to projects that have no policy collections assigned.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *iacPolicyCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.IaCPolicyCollectionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, diags := r.policies(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := iac.PolicyCollectionCreateInput{
		Name:      plan.Name.ValueString(),
		Policies:  policies,
		IsDefault: plan.IsDefault.ValueBool(),
	}
	collection, err := r.api.PolicyCollection.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, collection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iacPolicyCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.IaCPolicyCollectionResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collection, err := r.api.PolicyCollection.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(ctx, collection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *iacPolicyCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.IaCPolicyCollectionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, diags := r.policies(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := iac.PolicyCollectionUpdateInput{
		Policies:  policies,
		IsDefault: plan.IsDefault.ValueBool(),
	}
	collection, err := r.api.PolicyCollection.Update(ctx, plan.Name.ValueString(), input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, collection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iacPolicyCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.IaCPolicyCollectionResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.PolicyCollection.Delete(ctx, state.Name.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *iacPolicyCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"name"})
}

func (r *iacPolicyCollectionResource) policies(ctx context.Context, plan models.IaCPolicyCollectionResource) ([]string, diag.Diagnostics) {
	var policies []string
	diags := plan.Policies.ElementsAs(ctx, &policies, false)
	return policies, diags
}
