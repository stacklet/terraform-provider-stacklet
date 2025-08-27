// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ resource.Resource                = &configurationProfileTeamsResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileTeamsResource{}
	_ resource.ResourceWithImportState = &configurationProfileTeamsResource{}
)

type teamsWebhookSecret struct {
	Plaintext string
	Encrypted string
	Version   string
}

func NewConfigurationProfileTeamsResource() resource.Resource {
	return &configurationProfileTeamsResource{}
}

type configurationProfileTeamsResource struct {
	api *api.API
}

func (r *configurationProfileTeamsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_teams"
}

func (r *configurationProfileTeamsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the Microsoft Teams configuration profile.

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
		},
		Blocks: map[string]schema.Block{
			"webhook": schema.ListNestedBlock{
				Description: "Microsoft Teams webhook configuration.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The webhook name.",
							Required:    true,
						},
						"url": schema.StringAttribute{
							Description: "The encrypted webhook URL returned by the API.",
							Computed:    true,
						},
						"url_wo": schema.StringAttribute{
							Description: "The input value for the webhook URL.",
							Required:    true,
							Sensitive:   true,
							WriteOnly:   true,
						},
						"url_wo_version": schema.StringAttribute{
							Description: "The version for the webhook URL. Must be changed to update url_wo.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *configurationProfileTeamsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileTeamsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileTeamsResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	whSecretsState, diags := r.getWebhooksSecrets(ctx, state)
	resp.Diagnostics.Append(diags...)

	teamsConfig, err := r.api.ConfigurationProfile.ReadTeams(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateTeamsModel(&state, whSecretsState, teamsConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileTeamsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.ConfigurationProfileTeamsResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	whSecretsConfig, diags := r.getWebhooksSecrets(ctx, config)
	resp.Diagnostics.Append(diags...)

	webhooks, diags := r.getWebhooksForCreate(ctx, plan, whSecretsConfig)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := api.TeamsConfiguration{
		Webhooks: webhooks,
	}
	teamsConfig, err := r.api.ConfigurationProfile.UpsertTeams(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateTeamsModel(&plan, whSecretsConfig, teamsConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileTeamsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.ConfigurationProfileTeamsResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	whSecretsConfig, diags := r.getWebhooksSecrets(ctx, config)
	resp.Diagnostics.Append(diags...)
	whSecretsState, diags := r.getWebhooksSecrets(ctx, state)
	resp.Diagnostics.Append(diags...)

	webhooks, diags := r.getWebhooksForUpdate(ctx, plan, whSecretsConfig, whSecretsState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := api.TeamsConfiguration{
		Webhooks: webhooks,
	}
	teamsConfig, err := r.api.ConfigurationProfile.UpsertTeams(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateTeamsModel(&plan, whSecretsConfig, teamsConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileTeamsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileTeamsResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteTeams(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileTeamsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileTeams))...)
}

func (r configurationProfileTeamsResource) updateTeamsModel(m *models.ConfigurationProfileTeamsResource, webhookSecrets map[string]teamsWebhookSecret, config *api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(config.ID)
	m.Profile = types.StringValue(config.Profile)

	webhookVersions := map[string]string{}
	for name, secret := range webhookSecrets {
		webhookVersions[name] = secret.Version
	}

	// get the current names ordering to preserve it in the updated webhooks block
	names := models.ListItemsIdentifiers(m.Webhooks, "name")

	updater := modelupdate.NewConfigurationProfileUpdater(*config)
	webhooks, diags := updater.TeamsWebhooksWithSecret(webhookVersions, names)
	if diags.HasError() {
		return diags
	}
	m.Webhooks = webhooks
	return diags
}

func (r configurationProfileTeamsResource) getWebhooksSecrets(ctx context.Context, m models.ConfigurationProfileTeamsResource) (map[string]teamsWebhookSecret, diag.Diagnostics) {
	secrets := map[string]teamsWebhookSecret{}
	var diags diag.Diagnostics

	for i, elem := range m.Webhooks.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("webhook.%d", i)),
				"Invalid webhook configuration",
				"Webhook configuration block is invalid.",
			)
			return nil, diags
		}

		var w models.TeamsWebhookWithSecret
		if diags := block.As(ctx, &w, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		secrets[w.Name.ValueString()] = teamsWebhookSecret{
			Plaintext: w.URLWO.ValueString(),
			Encrypted: w.URL.ValueString(),
			Version:   w.URLWOVersion.ValueString(),
		}
	}
	return secrets, diags
}

func (r configurationProfileTeamsResource) getWebhooksForCreate(ctx context.Context, m models.ConfigurationProfileTeamsResource, webhookSecrets map[string]teamsWebhookSecret) ([]api.TeamsWebhook, diag.Diagnostics) {
	if m.Webhooks.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	webhooks := []api.TeamsWebhook{}
	for i, elem := range m.Webhooks.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("webhook.%d", i)),
				"Invalid webhook configuration",
				"Webhook configuration block is invalid.",
			)
			return nil, diags
		}
		var w models.TeamsWebhookWithSecret
		if diags := block.As(ctx, &w, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}

		webhooks = append(
			webhooks,
			api.TeamsWebhook{
				Name: w.Name.ValueString(),
				URL:  webhookSecrets[w.Name.ValueString()].Plaintext,
			},
		)
	}
	return webhooks, diags
}

func (r configurationProfileTeamsResource) getWebhooksForUpdate(ctx context.Context, m models.ConfigurationProfileTeamsResource, webhookSecretsFromConfig, webhookSecretsFromState map[string]teamsWebhookSecret) ([]api.TeamsWebhook, diag.Diagnostics) {
	if m.Webhooks.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	webhooks := []api.TeamsWebhook{}
	for i, elem := range m.Webhooks.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("webhook.%d", i)),
				"Invalid webhook configuration",
				"Webhook configuration block is invalid.",
			)
			return nil, diags
		}
		var w models.TeamsWebhookWithSecret
		if diags := block.As(ctx, &w, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}

		whConfig := webhookSecretsFromConfig[w.Name.ValueString()]
		whState := webhookSecretsFromState[w.Name.ValueString()]
		var url string
		if whState.Version == w.URLWOVersion.ValueString() {
			url = whState.Encrypted // send back the encrypted value
		} else {
			url = whConfig.Plaintext // sent the new value from the config
		}

		webhooks = append(
			webhooks,
			api.TeamsWebhook{
				Name: w.Name.ValueString(),
				URL:  url,
			},
		)
	}
	return webhooks, diags
}
