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

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &configurationProfileSlackResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileSlackResource{}
	_ resource.ResourceWithImportState = &configurationProfileSlackResource{}
)

type slackWebhookSecret struct {
	Plaintext string
	Encrypted string
	Version   string
}

func NewConfigurationProfileSlackResource() resource.Resource {
	return &configurationProfileSlackResource{}
}

type configurationProfileSlackResource struct {
	api *api.API
}

func (r *configurationProfileSlackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_slack"
}

func (r *configurationProfileSlackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the Slack configuration profile.

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
			"user_fields": schema.ListAttribute{
				Description: "Fields to use for identifying users for notification delivery.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     tftypes.EmptyListDefault(types.StringType),
			},
			"token": schema.StringAttribute{
				Description: "The encrypted value for the token, returned by the API.",
				Computed:    true,
			},
			"token_wo": schema.StringAttribute{
				Description: "The input value for the token.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"token_wo_version": schema.StringAttribute{
				Description: "The version for token. Must be changed to update token_wo.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"webhook": schema.ListNestedBlock{
				Description: "Slack webhook configuration.",
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

func (r *configurationProfileSlackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileSlackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileSlackResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	whSecretsState, diags := r.getWebhooksSecrets(ctx, state)
	resp.Diagnostics.Append(diags...)

	slackConfig, err := r.api.ConfigurationProfile.ReadSlack(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateSlackModel(ctx, &state, whSecretsState, slackConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileSlackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.ConfigurationProfileSlackResource
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

	input := api.SlackConfiguration{
		UserFields: api.StringsList(plan.UserFields),
		Token:      config.TokenWO.ValueStringPointer(),
		Webhooks:   webhooks,
	}
	slackConfig, err := r.api.ConfigurationProfile.UpsertSlack(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateSlackModel(ctx, &plan, whSecretsConfig, slackConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileSlackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.ConfigurationProfileSlackResource
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

	var token *string
	if state.TokenWOVersion == plan.TokenWOVersion {
		token = state.Token.ValueStringPointer() // send back the previous encrypted value (if set)
	} else {
		token = config.Token.ValueStringPointer() // send the new value from the config
	}

	input := api.SlackConfiguration{
		UserFields: api.StringsList(plan.UserFields),
		Token:      token,
		Webhooks:   webhooks,
	}
	slackConfig, err := r.api.ConfigurationProfile.UpsertSlack(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateSlackModel(ctx, &plan, whSecretsConfig, slackConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileSlackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileSlackResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteSlack(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileSlackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileSlack))...)
}

func (r configurationProfileSlackResource) updateSlackModel(ctx context.Context, m *models.ConfigurationProfileSlackResource, webhookSecrets map[string]slackWebhookSecret, cp *api.ConfigurationProfile) diag.Diagnostics {
	webhookVersions := map[string]string{}
	for name, secret := range webhookSecrets {
		webhookVersions[name] = secret.Version
	}

	return m.Update(ctx, *cp, webhookVersions)
}

func (r configurationProfileSlackResource) getWebhooksSecrets(ctx context.Context, m models.ConfigurationProfileSlackResource) (map[string]slackWebhookSecret, diag.Diagnostics) {
	secrets := map[string]slackWebhookSecret{}
	var diags diag.Diagnostics

	for i, elem := range m.Webhooks.Elements() {
		block, ok := elem.(types.Object)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("webhook.%d", i)),
				"Invalid webhook configuration",
				"Webhook configuration block is invalid.",
			)
			return nil, diags
		}

		var w models.SlackWebhookWithSecret
		if diags := block.As(ctx, &w, ObjectAsOptions); diags.HasError() {
			return nil, diags
		}
		secrets[w.Name.ValueString()] = slackWebhookSecret{
			Plaintext: w.URLWO.ValueString(),
			Encrypted: w.URL.ValueString(),
			Version:   w.URLWOVersion.ValueString(),
		}
	}
	return secrets, diags
}

func (r configurationProfileSlackResource) getWebhooksForCreate(ctx context.Context, m models.ConfigurationProfileSlackResource, webhookSecrets map[string]slackWebhookSecret) ([]api.SlackWebhook, diag.Diagnostics) {
	if m.Webhooks.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	webhooks := []api.SlackWebhook{}
	for i, elem := range m.Webhooks.Elements() {
		block, ok := elem.(types.Object)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("webhook.%d", i)),
				"Invalid webhook configuration",
				"Webhook configuration block is invalid.",
			)
			return nil, diags
		}
		var w models.SlackWebhookWithSecret
		if diags := block.As(ctx, &w, ObjectAsOptions); diags.HasError() {
			return nil, diags
		}

		webhooks = append(
			webhooks,
			api.SlackWebhook{
				Name: w.Name.ValueString(),
				URL:  webhookSecrets[w.Name.ValueString()].Plaintext,
			},
		)
	}
	return webhooks, diags
}

func (r configurationProfileSlackResource) getWebhooksForUpdate(ctx context.Context, m models.ConfigurationProfileSlackResource, webhookSecretsFromConfig, webhookSecretsFromState map[string]slackWebhookSecret) ([]api.SlackWebhook, diag.Diagnostics) {
	if m.Webhooks.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	webhooks := []api.SlackWebhook{}
	for i, elem := range m.Webhooks.Elements() {
		block, ok := elem.(types.Object)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("webhook.%d", i)),
				"Invalid webhook configuration",
				"Webhook configuration block is invalid.",
			)
			return nil, diags
		}
		var w models.SlackWebhookWithSecret
		if diags := block.As(ctx, &w, ObjectAsOptions); diags.HasError() {
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
			api.SlackWebhook{
				Name: w.Name.ValueString(),
				URL:  url,
			},
		)
	}
	return webhooks, diags
}
