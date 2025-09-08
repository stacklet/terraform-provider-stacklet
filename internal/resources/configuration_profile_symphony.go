// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ resource.Resource                = &configurationProfileSymphonyResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileSymphonyResource{}
	_ resource.ResourceWithImportState = &configurationProfileSymphonyResource{}
)

func NewConfigurationProfileSymphonyResource() resource.Resource {
	return &configurationProfileSymphonyResource{}
}

type configurationProfileSymphonyResource struct {
	api *api.API
}

func (r *configurationProfileSymphonyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_symphony"
}

func (r *configurationProfileSymphonyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the Symphony configuration profile.

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
			"agent_domain": schema.StringAttribute{
				Description: "The Symphony agent domain.",
				Required:    true,
			},
			"service_account": schema.StringAttribute{
				Description: "The Symphony service account.",
				Required:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The encrypted value for the account private key, returned from the API.",
				Computed:    true,
			},
			"private_key_wo": schema.StringAttribute{
				Description: "The input value for the account private key.",
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"private_key_wo_version": schema.StringAttribute{
				Description: "The version for the private key. Must be changed to update private_key_wo.",
				Required:    true,
			},
		},
	}
}

func (r *configurationProfileSymphonyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileSymphonyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileSymphonyResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.api.ConfigurationProfile.ReadSymphony(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	r.updateSymphonyModel(&state, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileSymphonyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.ConfigurationProfileSymphonyResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.SymphonyConfiguration{
		AgentDomain:    plan.AgentDomain.ValueString(),
		ServiceAccount: plan.ServiceAccount.ValueString(),
		PrivateKey:     config.PrivateKeyWO.ValueString(),
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertSymphony(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateSymphonyModel(&plan, profileConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileSymphonyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.ConfigurationProfileSymphonyResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var privateKey string
	if state.PrivateKeyWOVersion == plan.PrivateKeyWOVersion {
		privateKey = state.PrivateKey.ValueString() // send back the previous encrypted value
	} else {
		privateKey = config.PrivateKeyWO.ValueString() // send the new value from the config
	}

	input := api.SymphonyConfiguration{
		AgentDomain:    plan.AgentDomain.ValueString(),
		ServiceAccount: plan.ServiceAccount.ValueString(),
		PrivateKey:     privateKey,
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertSymphony(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateSymphonyModel(&plan, profileConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileSymphonyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileSymphonyResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteSymphony(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileSymphonyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileSymphony))...)
}

func (r configurationProfileSymphonyResource) updateSymphonyModel(m *models.ConfigurationProfileSymphonyResource, config *api.ConfigurationProfile) {
	symphonyConfig := config.Record.SymphonyConfiguration

	m.ID = types.StringValue(config.ID)
	m.Profile = types.StringValue(config.Profile)
	m.AgentDomain = types.StringValue(symphonyConfig.AgentDomain)
	m.ServiceAccount = types.StringValue(symphonyConfig.ServiceAccount)
	m.PrivateKey = types.StringValue(symphonyConfig.PrivateKey)
}
