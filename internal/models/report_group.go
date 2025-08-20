// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReportGroupResource is the model for notification report groups resources.
type ReportGroupResource struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Enabled                    types.Bool   `tfsdk:"enabled"`
	Bindings                   types.List   `tfsdk:"bindings"`
	Schedule                   types.String `tfsdk:"schedule"`
	GroupBy                    types.List   `tfsdk:"group_by"`
	UseMessageSettings         types.Bool   `tfsdk:"use_message_settings"`
	EmailDeliverySettings      types.List   `tfsdk:"email_delivery_settings"`
	SlackDeliverySettings      types.List   `tfsdk:"slack_delivery_settings"`
	TeamsDeliverySettings      types.List   `tfsdk:"teams_delivery_settings"`
	ServiceNowDeliverySettings types.List   `tfsdk:"servicenow_delivery_settings"`
	JiraDeliverySettings       types.List   `tfsdk:"jira_delivery_settings"`
	SymphonyDeliverySettings   types.List   `tfsdk:"symphony_delivery_settings"`
}

// ReportGroupDataSource is the model for notification report groups data sources.
type ReportGroupDataSource struct {
	ReportGroupResource

	Source types.String `tfsdk:"source"`
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

// TeamsDeliverySettings are the settings for a Teams delivery type.
type TeamsDeliverySettings struct {
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Recipients     types.List   `tfsdk:"recipients"`
	Template       types.String `tfsdk:"template"`
}

func (s TeamsDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only": types.BoolType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"template":         types.StringType,
	}
}

// ServiceNowDeliverySettings are the settings for a ServiceNow delivery type.
type ServiceNowDeliverySettings struct {
	FirstMatchOnly types.Bool   `tfsdk:"first_match_only"`
	Impact         types.String `tfsdk:"impact"`
	Recipients     types.List   `tfsdk:"recipients"`
	Description    types.String `tfsdk:"description"`
	Template       types.String `tfsdk:"template"`
	Urgency        types.String `tfsdk:"urgency"`
}

func (s ServiceNowDeliverySettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_match_only": types.BoolType,
		"impact":           types.StringType,
		"recipients":       types.ListType{ElemType: types.ObjectType{AttrTypes: Recipient{}.AttributeTypes()}},
		"description":      types.StringType,
		"template":         types.StringType,
		"urgency":          types.StringType,
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
