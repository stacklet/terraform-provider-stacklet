// Copyright (c) 2025 - Stacklet, Inc.

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
	"github.com/stacklet/terraform-provider-stacklet/internal/planmodifiers"
)

var (
	_ resource.Resource                = &notificationTemplateResource{}
	_ resource.ResourceWithConfigure   = &notificationTemplateResource{}
	_ resource.ResourceWithImportState = &notificationTemplateResource{}
)

func newNotificationTemplateResource() resource.Resource {
	return &notificationTemplateResource{}
}

type notificationTemplateResource struct {
	apiResource
}

func (r *notificationTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_template"
}

func (r *notificationTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a notification template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the notification template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the template.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the template.",
				Optional:    true,
			},
			"transport": schema.StringAttribute{
				Description: "The notification transport the template is for.",
				Optional:    true,
			},
			"content": schema.StringAttribute{
				Description: "The template content.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.TrimWhiteSpace(),
				},
			},
		},
	}
}

func (r *notificationTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.NotificationTemplateResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := r.api.Template.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(template)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *notificationTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.NotificationTemplateResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.TemplateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Transport:   plan.Transport.ValueStringPointer(),
		Content:     plan.Content.ValueString(),
	}
	template, err := r.api.Template.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(template)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *notificationTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.NotificationTemplateResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.TemplateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Transport:   plan.Transport.ValueStringPointer(),
		Content:     plan.Content.ValueString(),
	}
	template, err := r.api.Template.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(template)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *notificationTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.NotificationTemplateResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.Template.Delete(ctx, state.Name.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *notificationTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"name"})
}
