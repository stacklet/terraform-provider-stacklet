// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &datadogIntegrationResource{}
	_ resource.ResourceWithConfigure   = &datadogIntegrationResource{}
	_ resource.ResourceWithImportState = &datadogIntegrationResource{}
)

type datadogIntegrationResource struct {
	apiResource
}

func (r *datadogIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datadog_integration"
}

func (r *datadogIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the Datadog integration.

Supplying credentials makes the ` + "`datadog-metrics`" + ` policy filter available, so policies can select resources by metrics held in Datadog rather than by the cloud's own monitoring.

The integration applies to the whole deployment; adding multiple resources of this kind will cause them to override each other.

The credentials are write-only: they are sent to the API and never returned by it, so they aren't kept in Terraform state. ` + "`configured`" + ` reports whether the API holds them.

An import brings in the settings but not the key versions, for the same reason. The first apply after an import therefore sends whatever ` + "`api_key_wo`" + ` and ` + "`app_key_wo`" + ` hold, replacing the stored credentials.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A fixed identifier for the integration, which is global and has no ID of its own.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether policies may use the `datadog-metrics` filter. Defaults to `true`. Enabling requires credentials to be stored, or supplied in the same apply. Disabling keeps them, so the integration can be switched back on without re-entering them.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"site": schema.StringAttribute{
				Description: "The Datadog site the credentials belong to, as a bare hostname, e.g. `datadoghq.com`, `datadoghq.eu` or `us3.datadoghq.com`. Defaults to `datadoghq.com`. Not a URL: queries are sent to `https://api.<site>`. Must be one of [Datadog's own sites](https://docs.datadoghq.com/getting_started/site/); an unrecognized value is rejected when the resource is applied, rather than at plan time. Changing this rechecks the stored credentials against the new site, since a key is only valid for the site it was issued for.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("datadoghq.com"),
			},
			"configured": schema.BoolAttribute{
				Description: "Whether the API holds credentials. This is how to tell, since they are never returned. One flag rather than one per key, because the two are always stored together.",
				Computed:    true,
			},
			"status": models.StatusInfo{}.ResourceSchemaAttribute(
				"What the entry describes: `<region>-secret`, e.g. `us-east-1-secret`, is whether Secrets Manager holds the credentials in that region, which it replicates on its own schedule; `<region>-execution` is whether the region has recorded where to find them, so a policy running there can use them. Every region the configuration is meant to reach has both from the moment that is decided, so the list doesn't grow as propagation lands.",
			),
			// After this, write-only secrets and associated trigger attrs.
			"api_key_wo": schema.StringAttribute{
				Description: "The Datadog API key, which identifies the organization. Must be supplied together with `app_key_wo` the first time, since the metrics query the filter runs needs both.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"api_key_wo_version": schema.StringAttribute{
				Description: "Change value to update api_key_wo.",
				Optional:    true,
			},
			"app_key_wo": schema.StringAttribute{
				Description: "The Datadog application key, which carries the scopes authorizing the metrics query. Datadog only shows an application key when it is created, so keep it somewhere it can be read from, such as a secrets manager.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"app_key_wo_version": schema.StringAttribute{
				Description: "Change value to update app_key_wo.",
				Optional:    true,
			},
		},
	}
}

func (r *datadogIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.DatadogIntegrationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.DatadogIntegrationInput{
		Enabled: plan.Enabled.ValueBoolPointer(),
		Site:    plan.Site.ValueStringPointer(),
		APIKey:  config.APIKeyWO.ValueStringPointer(),
		AppKey:  config.AppKeyWO.ValueStringPointer(),
	}

	integration, err := r.api.DatadogIntegration.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *datadogIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.DatadogIntegrationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.api.DatadogIntegration.Read(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *datadogIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.DatadogIntegrationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only send a key whose version says it has changed: omitted fields keep
	// their stored value, which is the only way to change the site or the
	// enabled flag without resupplying credentials.
	input := api.DatadogIntegrationInput{
		Enabled: plan.Enabled.ValueBoolPointer(),
		Site:    plan.Site.ValueStringPointer(),
	}
	if !state.APIKeyWOVersion.Equal(plan.APIKeyWOVersion) {
		input.APIKey = config.APIKeyWO.ValueStringPointer()
	}
	if !state.AppKeyWOVersion.Equal(plan.AppKeyWOVersion) {
		input.AppKey = config.AppKeyWO.ValueStringPointer()
	}

	integration, err := r.api.DatadogIntegration.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *datadogIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DatadogIntegrationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.DatadogIntegration.Remove(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *datadogIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), models.DatadogIntegrationID)...)
}
