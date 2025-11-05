// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// ReportGroupDataSource is the model for notification report groups data sources.
type ReportGroupDataSource struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Enabled                    types.Bool   `tfsdk:"enabled"`
	Bindings                   types.List   `tfsdk:"bindings"`
	Schedule                   types.String `tfsdk:"schedule"`
	GroupBy                    types.List   `tfsdk:"group_by"`
	UseMessageSettings         types.Bool   `tfsdk:"use_message_settings"`
	EmailDeliverySettings      types.List   `tfsdk:"email_delivery_settings"`
	SlackDeliverySettings      types.List   `tfsdk:"slack_delivery_settings"`
	MSTeamsDeliverySettings    types.List   `tfsdk:"msteams_delivery_settings"`
	ServiceNowDeliverySettings types.List   `tfsdk:"servicenow_delivery_settings"`
	JiraDeliverySettings       types.List   `tfsdk:"jira_delivery_settings"`
	SymphonyDeliverySettings   types.List   `tfsdk:"symphony_delivery_settings"`
	Source                     types.String `tfsdk:"source"`
}

func (m *ReportGroupDataSource) Update(rg api.ReportGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(rg.ID)
	m.Name = types.StringValue(rg.Name)
	m.Enabled = types.BoolValue(rg.Enabled)
	m.Bindings = typehelpers.StringsList(rg.Bindings)
	m.Source = types.StringValue(string(rg.Source))
	m.Schedule = types.StringValue(rg.Schedule)
	m.GroupBy = typehelpers.StringsList(rg.GroupBy)
	m.UseMessageSettings = types.BoolValue(rg.UseMessageSettings)

	emailDeliverySettings, d := typehelpers.ObjectList[EmailDeliverySettings](
		rg.EmailDeliverySettings(),
		func(entry api.EmailDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := typehelpers.ObjectList[Recipient](
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
				"cc":               typehelpers.StringsList(entry.CC),
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
	diags.Append(d...)
	m.EmailDeliverySettings = emailDeliverySettings

	slackDeliverySettings, d := typehelpers.ObjectList[SlackDeliverySettings](
		rg.SlackDeliverySettings(),
		func(entry api.SlackDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := typehelpers.ObjectList[Recipient](
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
	diags.Append(d...)
	m.SlackDeliverySettings = slackDeliverySettings

	msteamsDeliverySettings, d := typehelpers.ObjectList[MSTeamsDeliverySettings](
		rg.MSTeamsDeliverySettings(),
		func(entry api.MSTeamsDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := typehelpers.ObjectList[Recipient](
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
	diags.Append(d...)
	m.MSTeamsDeliverySettings = msteamsDeliverySettings

	servicenowDeliverySettings, d := typehelpers.ObjectList[ServiceNowDeliverySettings](
		rg.ServiceNowDeliverySettings(),
		func(entry api.ServiceNowDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := typehelpers.ObjectList[Recipient](
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
	diags.Append(d...)
	m.ServiceNowDeliverySettings = servicenowDeliverySettings

	jiraDeliverySettings, d := typehelpers.ObjectList[JiraDeliverySettings](
		rg.JiraDeliverySettings(),
		func(entry api.JiraDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := typehelpers.ObjectList[Recipient](
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
	diags.Append(d...)
	m.JiraDeliverySettings = jiraDeliverySettings

	symphonyDeliverySettings, d := typehelpers.ObjectList[SymphonyDeliverySettings](
		rg.SymphonyDeliverySettings(),
		func(entry api.SymphonyDeliverySettings) (map[string]attr.Value, diag.Diagnostics) {
			recipients, diags := typehelpers.ObjectList[Recipient](
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
	diags.Append(d...)
	m.SymphonyDeliverySettings = symphonyDeliverySettings

	return diags
}

// ReportGroupResource is the model for notification report groups resources.
type ReportGroupResource struct {
	ReportGroupDataSource
}

// Recipient is the models for a notification recipient.
type Recipient struct {
	AccountOwner  types.Bool   `tfsdk:"account_owner"`
	EventOwner    types.Bool   `tfsdk:"event_owner"`
	ResourceOwner types.Bool   `tfsdk:"resource_owner"`
	Tag           types.String `tfsdk:"tag"`
	Value         types.String `tfsdk:"value"`
}

func (r Recipient) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_owner":  types.BoolType,
		"event_owner":    types.BoolType,
		"resource_owner": types.BoolType,
		"tag":            types.StringType,
		"value":          types.StringType,
	}
}

// EmailDeliverySettings are the settings for an email delivery type.
type EmailDeliverySettings struct {
	CC             types.List   `tfsdk:"cc"`
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Format         types.String `tfsdk:"format"`
	From           types.String `tfsdk:"from"`
	Priority       types.String `tfsdk:"priority"`
	Recipients     types.List   `tfsdk:"recipients"`
	Subject        types.String `tfsdk:"subject"`
	Template       types.String `tfsdk:"template"`
}

func (s EmailDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cc":               types.ListType{ElemType: types.StringType},
		"first_match_only": types.BoolType,
		"format":           types.StringType,
		"from":             types.StringType,
		"priority":         types.StringType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"subject":          types.StringType,
		"template":         types.StringType,
	}
}

// SlackDeliverySettings are the settings for a Slack delivery type.
type SlackDeliverySettings struct {
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Recipients     types.List   `tfsdk:"recipients"`
	Template       types.String `tfsdk:"template"`
}

func (s SlackDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only": types.BoolType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"template":         types.StringType,
	}
}

// MSTeamsDeliverySettings are the settings for a Microsoft Teams delivery type.
type MSTeamsDeliverySettings struct {
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Recipients     types.List   `tfsdk:"recipients"`
	Template       types.String `tfsdk:"template"`
}

func (s MSTeamsDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only": types.BoolType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"template":         types.StringType,
	}
}

// ServiceNowDeliverySettings are the settings for a ServiceNow delivery type.
type ServiceNowDeliverySettings struct {
	FirstMatchOnly   types.Bool   `tfsdk:"first_match_only"`
	Impact           types.String `tfsdk:"impact"`
	Recipients       types.List   `tfsdk:"recipients"`
	ShortDescription types.String `tfsdk:"short_description"`
	Template         types.String `tfsdk:"template"`
	Urgency          types.String `tfsdk:"urgency"`
}

func (s ServiceNowDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only":  types.BoolType,
		"impact":            types.StringType,
		"recipients":        types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"short_description": types.StringType,
		"template":          types.StringType,
		"urgency":           types.StringType,
	}
}

// JiraDeliverySettings are the settings for a Jira delivery type.
type JiraDeliverySettings struct {
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Recipients     types.List   `tfsdk:"recipients"`
	Template       types.String `tfsdk:"template"`
	Description    types.String `tfsdk:"description"`
	Project        types.String `tfsdk:"project"`
	Summary        types.String `tfsdk:"summary"`
}

func (s JiraDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only": types.BoolType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"template":         types.StringType,
		"description":      types.StringType,
		"project":          types.StringType,
		"summary":          types.StringType,
	}
}

// SymphonyDeliverySettings are the settings for a Symphony delivery type.
type SymphonyDeliverySettings struct {
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Recipients     types.List   `tfsdk:"recipients"`
	Template       types.String `tfsdk:"template"`
}

func (s SymphonyDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only": types.BoolType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"template":         types.StringType,
	}
}
