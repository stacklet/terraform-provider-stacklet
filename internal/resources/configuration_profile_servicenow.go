// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &configurationProfileServiceNowResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileServiceNowResource{}
	_ resource.ResourceWithImportState = &configurationProfileServiceNowResource{}
)

func newConfigurationProfileServiceNowResource() resource.Resource {
	return &configurationProfileServiceNowResource{}
}

type configurationProfileServiceNowResource struct {
	apiResource
}

func (r *configurationProfileServiceNowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_servicenow"
}

func (r *configurationProfileServiceNowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the ServiceNow configuration profile.

The profile is global, adding multiple resources of this kind will cause them to override each other.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint": schema.StringAttribute{
				Description: "The ServiceNow instance endpoint.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The ServiceNow instance authentication username.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The encrypted value for the ServiceNow instance authentication password, returned by the API.",
				Computed:    true,
			},
			"password_wo": schema.StringAttribute{
				Description: "The input value for the ServiceNow instance authentication password.",
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"password_wo_version": schema.StringAttribute{
				Description: "The version for the authentication password. Must be changed to update password_wo.",
				Required:    true,
			},
			"issue_type": schema.StringAttribute{
				Description: "The type of issue to use for tickets.",
				Required:    true,
			},
			"closed_state": schema.StringAttribute{
				Description: "The state for closed tickets.",
				Required:    true,
			},
		},
	}
}

func (r *configurationProfileServiceNowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileServiceNowResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileConfig, err := r.api.ConfigurationProfile.ReadServiceNow(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(*profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileServiceNowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.ConfigurationProfileServiceNowResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.ServiceNowConfiguration{
		Endpoint:    plan.Endpoint.ValueString(),
		User:        plan.Username.ValueString(),
		Password:    config.PasswordWO.ValueString(),
		IssueType:   plan.IssueType.ValueString(),
		ClosedState: plan.ClosedState.ValueString(),
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertServiceNow(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(*profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileServiceNowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.ConfigurationProfileServiceNowResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var password string
	if state.PasswordWOVersion == plan.PasswordWOVersion {
		password = state.Password.ValueString() // send back the previous encrypted value
	} else {
		password = config.PasswordWO.ValueString() // send the new value from the config
	}

	input := api.ServiceNowConfiguration{
		Endpoint:    plan.Endpoint.ValueString(),
		User:        plan.Username.ValueString(),
		Password:    password,
		IssueType:   plan.IssueType.ValueString(),
		ClosedState: plan.ClosedState.ValueString(),
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertServiceNow(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(*profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileServiceNowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileServiceNowResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteServiceNow(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileServiceNowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileServiceNow))...)
}
