// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &anthropicIntegrationResource{}
	_ resource.ResourceWithConfigure   = &anthropicIntegrationResource{}
	_ resource.ResourceWithImportState = &anthropicIntegrationResource{}
)

type anthropicIntegrationResource struct {
	apiResource
}

func (r *anthropicIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_anthropic_integration"
}

func (r *anthropicIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage an Anthropic integration, which imports the daily cost and usage of one Anthropic organization.

The admin key is checked against Anthropic when it is saved: one that doesn't authenticate, or belongs to another organization, fails the apply with the reason, and isn't stored.

There is one integration per organization. Creating this resource for an organization that already has an integration takes that integration over, replacing its admin key.

The admin key is write-only: it is sent to the API and never returned by it, so it isn't kept in Terraform state, and is sent only when admin_key_wo_version changes. The configured attribute reports whether the API holds one.

An import brings in the integration but not the key version, for the same reason. The first apply after an import therefore sends whatever admin_key_wo holds, replacing the stored key.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the Anthropic integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the Anthropic organization, as shown in the Claude Console's organization settings. The admin key must belong to it. Changing it replaces the integration.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"configured": schema.BoolAttribute{
				Description: "Whether the API holds an admin key. This is how to tell, since the key is never returned.",
				Computed:    true,
			},
			// After this, write-only secrets and associated trigger attrs.
			"admin_key_wo": schema.StringAttribute{
				Description: "An Anthropic Admin API key for the organization. Anthropic only shows an admin key when it is created, so keep it somewhere it can be read from, such as a secrets manager.",
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"admin_key_wo_version": schema.StringAttribute{
				Description: "The version for the admin key. Must be changed to update admin_key_wo.",
				Required:    true,
			},
		},
	}
}

func (r *anthropicIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.AnthropicIntegrationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AnthropicIntegrationInput{
		OrgID:    plan.OrgID.ValueString(),
		AdminKey: config.AdminKeyWO.ValueStringPointer(),
	}
	integration, err := r.api.AnthropicIntegration.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *anthropicIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AnthropicIntegrationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.api.AnthropicIntegration.Read(ctx, state.OrgID.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *anthropicIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.AnthropicIntegrationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The organization can't change in place, so the only update is a new key,
	// sent when its version says there is one.
	input := api.AnthropicIntegrationInput{
		OrgID:    plan.OrgID.ValueString(),
		AdminKey: credentialToSend(&resp.Diagnostics, !state.AdminKeyWOVersion.Equal(plan.AdminKeyWOVersion), "admin_key_wo", config.AdminKeyWO),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.api.AnthropicIntegration.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *anthropicIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AnthropicIntegrationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.AnthropicIntegration.Remove(ctx, state.OrgID.ValueString())
	if conflict, ok := err.(api.Conflict); ok {
		// The integration is gone, but its stored key outlived it. Failing the
		// destroy would keep a resource the next refresh drops anyway, so
		// Terraform could never retry it; the API finishes the cleanup the next
		// time an integration for this organization is saved or removed.
		resp.Diagnostics.AddWarning("Anthropic integration removed with its admin key pending deletion", conflict.Error())
		return
	}
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	}
}

func (r *anthropicIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"org_id"})
}
