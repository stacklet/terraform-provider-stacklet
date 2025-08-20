// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &reportGroupResource{}
	_ resource.ResourceWithConfigure   = &reportGroupResource{}
	_ resource.ResourceWithImportState = &reportGroupResource{}
)

func NewReportGroupResource() resource.Resource {
	return &reportGroupResource{}
}

type reportGroupResource struct {
	api *api.API
}

func (r *reportGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_report_group"
}

func (r *reportGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a notification report group by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the report group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name for the report group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the report group is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"bindings": schema.ListAttribute{
				Description: "List of UUIDs for bindings the report group is for.",
				Required:    true,
				ElementType: types.StringType,
			},
			"schedule": schema.StringAttribute{
				Description: "Notification schedule.",
				Required:    true,
			},
			"group_by": schema.ListAttribute{
				Description: "Fields on which matching resources are grouped.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(basetypes.NewListValueMust(types.StringType, []attr.Value{})),
			},
			"use_message_settings": schema.BoolAttribute{
				Description: "Whether to use delivery settings from the notification message.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
		Blocks: map[string]schema.Block{
			"email_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for email.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cc": schema.ListAttribute{
							Description: "List of CC addresses.",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Default:     listdefault.StaticValue(basetypes.NewListValueMust(types.StringType, []attr.Value{})),
						},
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"format": schema.StringAttribute{
							Description: "Email format (html or plain). Autodetected from the template if not specified.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"from": schema.StringAttribute{
							Description: "Email from address.",
							Optional:    true,
						},
						"priority": schema.StringAttribute{
							Description: "Email priority.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("3"),
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Required:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Optional:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Optional:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Optional:    true,
									},
								},
								Validators: []validator.Object{
									validRecipient{},
								},
							},
						},
						"subject": schema.StringAttribute{
							Description: "Email subject.",
							Required:    true,
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the email.",
							Required:    true,
						},
					},
				},
			},
			"slack_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for Slack.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Required:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Optional:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Optional:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Optional:    true,
									},
								},
								Validators: []validator.Object{
									validRecipient{},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Required:    true,
						},
					},
				},
			},
			"teams_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for Microsoft Teams.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Optional:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Optional:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Optional:    true,
									},
								},
								Validators: []validator.Object{
									validRecipient{},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Required:    true,
						},
					},
				},
			},
			"servicenow_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for ServiceNow.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Optional:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Optional:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Optional:    true,
									},
								},
								Validators: []validator.Object{
									validRecipient{},
								},
							},
						},
						"impact": schema.StringAttribute{
							Description: "Impact to use for the ticket.",
							Required:    true,
						},
						"short_description": schema.StringAttribute{
							Description: "Ticket description.",
							Required:    true,
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Required:    true,
						},
						"urgency": schema.StringAttribute{
							Description: "Ticket urgency.",
							Required:    true,
						},
					},
				},
			},
			"jira_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for Jira.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Optional:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Optional:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Optional:    true,
									},
								},
								Validators: []validator.Object{
									validRecipient{},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "Ticket description.",
							Required:    true,
						},
						"project": schema.StringAttribute{
							Description: "Jira project key.",
							Required:    true,
						},
						"summary": schema.StringAttribute{
							Description: "Ticket summary.",
							Required:    true,
						},
					},
				},
			},
			"symphony_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for Symphony.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Optional:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Optional:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Optional:    true,
									},
								},
								Validators: []validator.Object{
									validRecipient{},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *reportGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *reportGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ReportGroupResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reportGroup, err := r.api.ReportGroup.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateReportGroupModel(&state, reportGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *reportGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ReportGroupResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailSettings, diags := r.getEmailDeliverySettings(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.ReportGroupInput{
		Name:               plan.Name.ValueString(),
		Enabled:            plan.Enabled.ValueBool(),
		Bindings:           api.StringsList(plan.Bindings),
		Source:             api.ReportSourceBinding, // only binding-based report groups are managed for now
		Schedule:           plan.Schedule.ValueString(),
		GroupBy:            api.StringsList(plan.GroupBy),
		UseMessageSettings: plan.UseMessageSettings.ValueBool(),
		EmailSettings:      emailSettings,
	}
	reportGroup, err := r.api.ReportGroup.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateReportGroupModel(&plan, reportGroup)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *reportGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ReportGroupResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailSettings, diags := r.getEmailDeliverySettings(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.ReportGroupInput{
		Name:               plan.Name.ValueString(),
		Enabled:            plan.Enabled.ValueBool(),
		Bindings:           api.StringsList(plan.Bindings),
		Source:             api.ReportSourceBinding, // only binding-based report groups are managed for now
		Schedule:           plan.Schedule.ValueString(),
		GroupBy:            api.StringsList(plan.GroupBy),
		UseMessageSettings: plan.UseMessageSettings.ValueBool(),
		EmailSettings:      emailSettings,
	}
	reportGroup, err := r.api.ReportGroup.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateReportGroupModel(&plan, reportGroup)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *reportGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ReportGroupResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ReportGroup.Delete(ctx, state.Name.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *reportGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

func (r reportGroupResource) updateReportGroupModel(m *models.ReportGroupResource, reportGroup *api.ReportGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(reportGroup.ID)
	m.Name = types.StringValue(reportGroup.Name)
	m.Enabled = types.BoolValue(reportGroup.Enabled)
	m.Bindings = tftypes.StringsList(reportGroup.Bindings)
	m.Schedule = types.StringValue(reportGroup.Schedule)
	m.GroupBy = tftypes.StringsList(reportGroup.GroupBy)
	m.UseMessageSettings = types.BoolValue(reportGroup.UseMessageSettings)

	updater := modelupdate.NewReportGroupUpdater(*reportGroup)

	emailDeliverySettings, d := updater.EmailDeliverySettings()
	diags.Append(d...)
	m.EmailDeliverySettings = emailDeliverySettings

	slackDeliverySettings, d := updater.SlackDeliverySettings()
	diags.Append(d...)
	m.SlackDeliverySettings = slackDeliverySettings

	teamsDeliverySettings, d := updater.TeamsDeliverySettings()
	diags.Append(d...)
	m.TeamsDeliverySettings = teamsDeliverySettings

	servicenowDeliverySettings, d := updater.ServiceNowDeliverySettings()
	diags.Append(d...)
	m.ServiceNowDeliverySettings = servicenowDeliverySettings

	jiraDeliverySettings, d := updater.JiraDeliverySettings()
	diags.Append(d...)
	m.JiraDeliverySettings = jiraDeliverySettings

	symphonyDeliverySettings, d := updater.SymphonyDeliverySettings()
	diags.Append(d...)
	m.SymphonyDeliverySettings = symphonyDeliverySettings

	return diags
}

func (r reportGroupResource) getEmailDeliverySettings(ctx context.Context, m models.ReportGroupResource) ([]api.EmailDeliverySettings, diag.Diagnostics) {
	if m.EmailDeliverySettings.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	settings := []api.EmailDeliverySettings{}
	for i, elem := range m.EmailDeliverySettings.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("email_delivery_settings.%d", i)),
				"Invalid email delivery settings",
				"Email delivery settings block is invalid.",
			)
			return nil, diags
		}
		var s models.EmailDeliverySettings
		if diags := block.As(ctx, &s, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}

		recipients, diags := r.getRecipients(ctx, s.Recipients)
		if diags.HasError() {
			return nil, diags
		}

		settings = append(
			settings,
			api.EmailDeliverySettings{
				CC:             api.StringsList(s.CC),
				FirstMatchOnly: s.FirstMatchOnly.ValueBoolPointer(),
				Format:         s.Format.ValueStringPointer(),
				FromEmail:      s.From.ValueStringPointer(),
				Priority:       s.Priority.ValueStringPointer(),
				Recipients:     recipients,
				Subject:        s.Subject.ValueString(),
				Template:       s.Template.ValueString(),
			},
		)
	}
	return settings, diags
}

func (r reportGroupResource) getRecipients(ctx context.Context, l types.List) ([]api.Recipient, diag.Diagnostics) {
	if l.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	recipients := []api.Recipient{}
	for i, elem := range l.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("recipients.%d", i)),
				"Invalid recipient",
				"Recipient block is invalid.",
			)
			return nil, diags
		}
		var r models.Recipient
		if diags := block.As(ctx, &r, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}

		recipients = append(
			recipients,
			api.Recipient{
				AccountOwner:  r.AccountOwner.ValueBoolPointer(),
				EventOwner:    r.EventOwner.ValueBoolPointer(),
				ResourceOwner: r.ResourceOwner.ValueBoolPointer(),
				Tag:           r.Tag.ValueStringPointer(),
				Value:         r.Value.ValueStringPointer(),
			},
		)
	}
	return recipients, nil
}

type validRecipient struct{}

func (v validRecipient) Description(ctx context.Context) string {
	return "Ensures exactly one recipient field is set: either one of the owners flags are true, or exactly one of tag and value is set"
}

func (v validRecipient) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validRecipient) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	setCount := 0
	for _, attr := range req.ConfigValue.Attributes() {
		if attr.IsNull() || attr.IsUnknown() {
			continue
		}
		switch a := attr.(type) {
		case types.Bool:
			if a.ValueBool() {
				setCount++
			}
		case types.String:
			if a.ValueString() != "" {
				setCount++
			}
		}
	}

	if setCount != 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid recipient configuration",
			"Exactly one recipient field must be set.",
		)
		return
	}
}
