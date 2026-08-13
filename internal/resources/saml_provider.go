// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &samlProviderResource{}
	_ resource.ResourceWithConfigure   = &samlProviderResource{}
	_ resource.ResourceWithImportState = &samlProviderResource{}
)

type samlProviderResource struct {
	apiResource
}

func (r *samlProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saml_provider"
}

func (r *samlProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a SAML identity provider, used to federate user login with an external identity provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the SAML provider.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name identifying the provider. It can't be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the SAML provider.",
				Optional:    true,
			},
			"metadata_url": schema.StringAttribute{
				Description: "URL of the SAML metadata endpoint. Exactly one of `metadata_url` or `metadata_xml` must be set.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("metadata_xml")),
				},
			},
			"metadata_xml": schema.StringAttribute{
				Description: "Inline SAML metadata XML document. Exactly one of `metadata_url` or `metadata_xml` must be set.",
				Optional:    true,
			},
			"enable_signout": schema.BoolAttribute{
				Description: "Whether the identity provider signout flow is enabled for this provider.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"idp_alias": schema.StringAttribute{
				Description: "A unique human-facing alias for the provider.",
				Optional:    true,
			},
		},
	}
}

func (r *samlProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SAMLProviderResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.SAMLProviderCreateInput{
		Name:          plan.Name.ValueString(),
		DisplayName:   plan.DisplayName.ValueStringPointer(),
		MetadataURL:   plan.MetadataURL.ValueStringPointer(),
		MetadataXML:   plan.MetadataXML.ValueStringPointer(),
		EnableSignout: plan.EnableSignout.ValueBool(),
		IDPAlias:      plan.IDPAlias.ValueStringPointer(),
	}

	samlProvider, err := r.api.SAMLProvider.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(samlProvider)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *samlProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.SAMLProviderResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	samlProvider, err := r.api.SAMLProvider.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(samlProvider)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *samlProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.SAMLProviderResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.SAMLProviderUpdateInput{
		Name:          plan.Name.ValueString(),
		DisplayName:   plan.DisplayName.ValueStringPointer(),
		MetadataURL:   plan.MetadataURL.ValueStringPointer(),
		MetadataXML:   plan.MetadataXML.ValueStringPointer(),
		EnableSignout: plan.EnableSignout.ValueBool(),
		IDPAlias:      plan.IDPAlias.ValueStringPointer(),
	}

	samlProvider, err := r.api.SAMLProvider.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(samlProvider)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *samlProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SAMLProviderResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.SAMLProvider.Delete(ctx, state.Name.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *samlProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"name"})
}
