// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ resource.Resource                = &configurationProfileEmailResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileEmailResource{}
	_ resource.ResourceWithImportState = &configurationProfileEmailResource{}
)

func NewConfigurationProfileEmailResource() resource.Resource {
	return &configurationProfileEmailResource{}
}

type configurationProfileEmailResource struct {
	api *api.API
}

func (r *configurationProfileEmailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_email"
}

func (r *configurationProfileEmailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the email configuration profile.

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
			"from": schema.StringAttribute{
				Description: "The email from field value.",
				Required:    true,
			},
			"ses_region": schema.StringAttribute{
				Description: "AWS SES region to use for sending emails.",
				Optional:    true,
			},
			"smtp": schema.SingleNestedAttribute{
				Description: "SMTP configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"server": schema.StringAttribute{
						Description: "SMTP server hostname or IP address.",
						Required:    true,
					},
					"port": schema.StringAttribute{
						Description: "SMTP server port.",
						Required:    true,
					},
					"ssl": schema.BoolAttribute{
						Description: "Whether SSL/TLS is enabled.",
						Optional:    true,
					},
					"username": schema.StringAttribute{
						Description: "Authentication username.",
						Optional:    true,
					},
					"password": schema.StringAttribute{
						Description: "Authentication password (encrypted), returned by the API.",
						Computed:    true,
					},
					"password_wo": schema.StringAttribute{
						Description: "Input value for authentication password.",
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
					},
					"password_wo_version": schema.StringAttribute{
						Description: "The version for the authentication password. Must be changed to update password_wo.",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r *configurationProfileEmailResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.ConfigurationProfileEmailResource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileConfig, err := r.api.ConfigurationProfile.ReadEmail(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, *profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configurationProfileEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.ConfigurationProfileEmailResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	smtpPlan, diags := r.getSMTPResource(ctx, plan)
	resp.Diagnostics.Append(diags...)
	smtpConfig, diags := r.getSMTPResource(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var smtp *api.SMTPConfiguration
	if !smtpPlan.Server.IsNull() {
		smtp = &api.SMTPConfiguration{
			Server:   smtpPlan.Server.ValueString(),
			Port:     smtpPlan.Port.ValueString(),
			SSL:      smtpPlan.SSL.ValueBoolPointer(),
			Username: smtpPlan.Username.ValueStringPointer(),
			Password: smtpConfig.PasswordWO.ValueStringPointer(),
		}
	}

	input := api.EmailConfiguration{
		FromEmail: plan.From.ValueString(),
		SESRegion: plan.SESRegion.ValueStringPointer(),
		SMTP:      smtp,
	}

	profileConfig, err := r.api.ConfigurationProfile.UpsertEmail(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, *profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.ConfigurationProfileEmailResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	smtpPlan, diags := r.getSMTPResource(ctx, plan)
	resp.Diagnostics.Append(diags...)
	smtpState, diags := r.getSMTPResource(ctx, state)
	resp.Diagnostics.Append(diags...)
	smtpConfig, diags := r.getSMTPResource(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var smtp *api.SMTPConfiguration
	if !smtpPlan.Server.IsNull() {

		var password *string
		if smtpState.PasswordWOVersion == smtpPlan.PasswordWOVersion {
			password = smtpState.Password.ValueStringPointer() // send back the previous encrypted value
		} else {
			password = smtpConfig.PasswordWO.ValueStringPointer() // send the new value from the config
		}

		smtp = &api.SMTPConfiguration{
			Server:   smtpPlan.Server.ValueString(),
			Port:     smtpPlan.Port.ValueString(),
			SSL:      smtpPlan.SSL.ValueBoolPointer(),
			Username: smtpPlan.Username.ValueStringPointer(),
			Password: password,
		}
	}

	input := api.EmailConfiguration{
		FromEmail: plan.From.ValueString(),
		SESRegion: plan.SESRegion.ValueStringPointer(),
		SMTP:      smtp,
	}

	profileConfig, err := r.api.ConfigurationProfile.UpsertEmail(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, *profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if err := r.api.ConfigurationProfile.DeleteEmail(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	}
}

func (r *configurationProfileEmailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileEmail))...)
}

func (r configurationProfileEmailResource) getSMTPResource(ctx context.Context, m models.ConfigurationProfileEmailResource) (models.SMTPResource, diag.Diagnostics) {
	var diags diag.Diagnostics
	var smtp models.SMTPResource

	if !m.SMTP.IsNull() {
		diags = m.SMTP.As(ctx, &smtp, basetypes.ObjectAsOptions{})
	}
	return smtp, diags
}
