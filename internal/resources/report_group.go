// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemavalidate"
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
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
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
							Optional:    true,
							ElementType: types.StringType,
						},
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Optional:    true,
						},
						"format": schema.StringAttribute{
							Description: "Email format (html or plain). Autodetected from the template if not specified.",
							Optional:    true,
						},
						"from": schema.StringAttribute{
							Description: "Email from address.",
							Optional:    true,
						},
						"priority": schema.StringAttribute{
							Description: "Email priority.",
							Optional:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
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
									schemavalidate.ExactlyOneRecipient(),
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
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
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
									schemavalidate.ExactlyOneRecipient(),
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
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
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
									schemavalidate.ExactlyOneRecipient(),
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
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
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
									schemavalidate.ExactlyOneRecipient(),
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
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
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
									schemavalidate.ExactlyOneRecipient(),
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
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Optional:    true,
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
									schemavalidate.ExactlyOneRecipient(),
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

	r.updateReportGroupModel(&state, reportGroup)
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

	input := api.ReportGroupInput{
		Name:               plan.Name.ValueString(),
		Enabled:            plan.Enabled.ValueBool(),
		Bindings:           api.StringsList(plan.Bindings),
		Source:             api.ReportSourceBinding, // only binding-based report groups are managed for now
		Schedule:           plan.Schedule.ValueString(),
		GroupBy:            api.StringsList(plan.GroupBy),
		UseMessageSettings: plan.UseMessageSettings.ValueBool(),
	}
	reportGroup, err := r.api.ReportGroup.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateReportGroupModel(&plan, reportGroup)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *reportGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ReportGroupResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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
	}
	reportGroup, err := r.api.ReportGroup.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateReportGroupModel(&plan, reportGroup)
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

func (r reportGroupResource) updateReportGroupModel(m *models.ReportGroupResource, reportGroup *api.ReportGroup) {
	m.ID = types.StringValue(reportGroup.ID)
	m.Name = types.StringValue(reportGroup.Name)
	m.Enabled = types.BoolValue(reportGroup.Enabled)
	m.Bindings = tftypes.StringsList(reportGroup.Bindings)
	m.Schedule = types.StringValue(reportGroup.Schedule)
	m.GroupBy = tftypes.StringsList(reportGroup.GroupBy)
	m.UseMessageSettings = types.BoolValue(reportGroup.UseMessageSettings)
}
