// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &reportGroupDataSource{}
)

func NewReportgroupDataSource() datasource.DataSource {
	return &reportGroupDataSource{}
}

type reportGroupDataSource struct {
	api *api.API
}

func (d *reportGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_report_group"
}

func (d *reportGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a notification report group by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the report group.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name for the report group.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the report group is enabled.",
				Computed:    true,
			},
			"bindings": schema.ListAttribute{
				Description: "List of UUIDs for bindings the report group is for.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source": schema.StringAttribute{
				Description: "Type of the source for the report group.",
				Computed:    true,
			},
			"schedule": schema.StringAttribute{
				Description: "Notification schedule.",
				Computed:    true,
			},
			"group_by": schema.ListAttribute{
				Description: "Fields on which matching resources are grouped.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"use_message_settings": schema.BoolAttribute{
				Description: "Whether to use delivery settings from the notification message.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"email_delivery_settings": schema.ListNestedBlock{
				Description: "Notifications delivery settings for email.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cc": schema.ListAttribute{
							Description: "List of CC addresses.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"first_match_only": schema.BoolAttribute{
							Description: "Only report the first match.",
							Computed:    true,
						},
						"format": schema.StringAttribute{
							Description: "Email format (html or plain). Autodetected from the template if not specified.",
							Computed:    true,
						},
						"from": schema.StringAttribute{
							Description: "Email from address.",
							Computed:    true,
						},
						"priority": schema.StringAttribute{
							Description: "Email priority.",
							Computed:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Computed:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Computed:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Computed:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Computed:    true,
									},
								},
							},
						},
						"subject": schema.StringAttribute{
							Description: "Email subject.",
							Computed:    true,
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the email.",
							Computed:    true,
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
							Computed:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Computed:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Computed:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Computed:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Computed:    true,
									},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Computed:    true,
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
							Computed:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Computed:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Computed:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Computed:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Computed:    true,
									},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Computed:    true,
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
							Computed:    true,
						},
						"impact": schema.StringAttribute{
							Description: "Impact to use for the ticket.",
							Computed:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Computed:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Computed:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Computed:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Computed:    true,
									},
								},
							},
						},
						"short_description": schema.StringAttribute{
							Description: "Ticket description.",
							Computed:    true,
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Computed:    true,
						},
						"urgency": schema.StringAttribute{
							Description: "Ticket urgency.",
							Computed:    true,
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
							Computed:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Computed:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Computed:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Computed:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Computed:    true,
									},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Ticket description.",
							Computed:    true,
						},
						"project": schema.StringAttribute{
							Description: "Jira project key.",
							Computed:    true,
						},
						"summary": schema.StringAttribute{
							Description: "Ticket summary.",
							Computed:    true,
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
							Computed:    true,
						},
						"recipients": schema.ListNestedAttribute{
							Description: "Recipients for the notification.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_owner": schema.BoolAttribute{
										Description: "Whether to notify the account owner.",
										Computed:    true,
									},
									"event_owner": schema.BoolAttribute{
										Description: "Whether to notify the event owner.",
										Computed:    true,
									},
									"resource_owner": schema.BoolAttribute{
										Description: "Whether to notify the resource owner.",
										Computed:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag to match the resource owner from.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Explicit value for a notification recipient.",
										Computed:    true,
									},
								},
							},
						},
						"template": schema.StringAttribute{
							Description: "Name of the template for the notification.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *reportGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *reportGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ReportGroupDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reportGroup, err := d.api.ReportGroup.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(reportGroup.ID)
	data.Name = types.StringValue(reportGroup.Name)
	data.Enabled = types.BoolValue(reportGroup.Enabled)
	data.Bindings = tftypes.StringsList(reportGroup.Bindings)
	data.Source = types.StringValue(string(reportGroup.Source))
	data.Schedule = types.StringValue(reportGroup.Schedule)
	data.GroupBy = tftypes.StringsList(reportGroup.GroupBy)
	data.UseMessageSettings = types.BoolValue(reportGroup.UseMessageSettings)

	emailDeliverySettings, diags := tftypes.ObjectList[models.EmailDeliverySettings](
		reportGroup.EmailDeliverySettings(),
		func(entry api.EmailDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := tftypes.ObjectList[models.Recipient](
				entry.Recipients,
				func(entry api.Recipient) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"account_owner":  types.BoolPointerValue(entry.AccountOwner),
						"event_owner":    types.BoolPointerValue(entry.EventOwner),
						"resource_owner": types.BoolPointerValue(entry.ResourceOwner),
						"tag":            types.StringPointerValue(entry.Tag),
						"value":          types.StringPointerValue(entry.Value),
					}, nil
				},
			)
			if diags.HasError() {
				return map[string]attr.Value{}, diags
			}

			return map[string]attr.Value{
				"cc":               tftypes.StringsList(entry.CC),
				"first_match_only": types.BoolPointerValue(entry.FirstMatchOnly),
				"format":           types.StringPointerValue(entry.Format),
				"from":             types.StringPointerValue(entry.FromEmail),
				"priority":         types.StringPointerValue(entry.Priority),
				"recipients":       recipients,
				"subject":          types.StringValue(entry.Subject),
				"template":         types.StringValue(entry.Template),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.EmailDeliverySettings = emailDeliverySettings

	slackDeliverySettings, diags := tftypes.ObjectList[models.SlackDeliverySettings](
		reportGroup.SlackDeliverySettings(),
		func(entry api.SlackDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := tftypes.ObjectList[models.Recipient](
				entry.Recipients,
				func(entry api.Recipient) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"account_owner":  types.BoolPointerValue(entry.AccountOwner),
						"event_owner":    types.BoolPointerValue(entry.EventOwner),
						"resource_owner": types.BoolPointerValue(entry.ResourceOwner),
						"tag":            types.StringPointerValue(entry.Tag),
						"value":          types.StringPointerValue(entry.Value),
					}, nil
				},
			)
			if diags.HasError() {
				return map[string]attr.Value{}, diags
			}

			return map[string]attr.Value{
				"first_match_only": types.BoolPointerValue(entry.FirstMatchOnly),
				"recipients":       recipients,
				"template":         types.StringValue(entry.Template),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.SlackDeliverySettings = slackDeliverySettings

	teamsDeliverySettings, diags := tftypes.ObjectList[models.TeamsDeliverySettings](
		reportGroup.TeamsDeliverySettings(),
		func(entry api.TeamsDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := tftypes.ObjectList[models.Recipient](
				entry.Recipients,
				func(entry api.Recipient) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"account_owner":  types.BoolPointerValue(entry.AccountOwner),
						"event_owner":    types.BoolPointerValue(entry.EventOwner),
						"resource_owner": types.BoolPointerValue(entry.ResourceOwner),
						"tag":            types.StringPointerValue(entry.Tag),
						"value":          types.StringPointerValue(entry.Value),
					}, nil
				},
			)
			if diags.HasError() {
				return map[string]attr.Value{}, diags
			}

			return map[string]attr.Value{
				"first_match_only": types.BoolPointerValue(entry.FirstMatchOnly),
				"recipients":       recipients,
				"template":         types.StringValue(entry.Template),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.TeamsDeliverySettings = teamsDeliverySettings

	servicenowDeliverySettings, diags := tftypes.ObjectList[models.ServiceNowDeliverySettings](
		reportGroup.ServiceNowDeliverySettings(),
		func(entry api.ServiceNowDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := tftypes.ObjectList[models.Recipient](
				entry.Recipients,
				func(entry api.Recipient) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"account_owner":  types.BoolPointerValue(entry.AccountOwner),
						"event_owner":    types.BoolPointerValue(entry.EventOwner),
						"resource_owner": types.BoolPointerValue(entry.ResourceOwner),
						"tag":            types.StringPointerValue(entry.Tag),
						"value":          types.StringPointerValue(entry.Value),
					}, nil
				},
			)
			if diags.HasError() {
				return map[string]attr.Value{}, diags
			}

			return map[string]attr.Value{
				"first_match_only":  types.BoolPointerValue(entry.FirstMatchOnly),
				"impact":            types.StringValue(entry.Impact),
				"recipients":        recipients,
				"short_description": types.StringValue(entry.ShortDescription),
				"template":          types.StringValue(entry.Template),
				"urgency":           types.StringValue(entry.Urgency),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ServiceNowDeliverySettings = servicenowDeliverySettings

	jiraDeliverySettings, diags := tftypes.ObjectList[models.JiraDeliverySettings](
		reportGroup.JiraDeliverySettings(),
		func(entry api.JiraDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := tftypes.ObjectList[models.Recipient](
				entry.Recipients,
				func(entry api.Recipient) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"account_owner":  types.BoolPointerValue(entry.AccountOwner),
						"event_owner":    types.BoolPointerValue(entry.EventOwner),
						"resource_owner": types.BoolPointerValue(entry.ResourceOwner),
						"tag":            types.StringPointerValue(entry.Tag),
						"value":          types.StringPointerValue(entry.Value),
					}, nil
				},
			)
			if diags.HasError() {
				return map[string]attr.Value{}, diags
			}

			return map[string]attr.Value{
				"first_match_only": types.BoolPointerValue(entry.FirstMatchOnly),
				"recipients":       recipients,
				"template":         types.StringValue(entry.Template),
				"description":      types.StringValue(entry.Description),
				"project":          types.StringValue(entry.Project),
				"summary":          types.StringValue(entry.Summary),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.JiraDeliverySettings = jiraDeliverySettings

	symphonyDeliverySettings, diags := tftypes.ObjectList[models.SymphonyDeliverySettings](
		reportGroup.SymphonyDeliverySettings(),
		func(entry api.SymphonyDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := tftypes.ObjectList[models.Recipient](
				entry.Recipients,
				func(entry api.Recipient) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"account_owner":  types.BoolPointerValue(entry.AccountOwner),
						"event_owner":    types.BoolPointerValue(entry.EventOwner),
						"resource_owner": types.BoolPointerValue(entry.ResourceOwner),
						"tag":            types.StringPointerValue(entry.Tag),
						"value":          types.StringPointerValue(entry.Value),
					}, nil
				},
			)
			if diags.HasError() {
				return map[string]attr.Value{}, diags
			}

			return map[string]attr.Value{
				"first_match_only": types.BoolPointerValue(entry.FirstMatchOnly),
				"recipients":       recipients,
				"template":         types.StringValue(entry.Template),
			}, nil
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.SymphonyDeliverySettings = symphonyDeliverySettings

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
