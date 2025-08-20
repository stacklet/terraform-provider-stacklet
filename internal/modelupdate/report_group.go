// Copyright (c) 2025 - Stacklet, Inc.

package modelupdate

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

// NewReportGroupUpdater returns a report group updater helper.
func NewReportGroupUpdater(rg api.ReportGroup) reportGroupUpdater {
	return reportGroupUpdater{rg: rg}
}

type reportGroupUpdater struct {
	rg api.ReportGroup
}

// EmailDeliverySettings returns email delivery settings.
func (u reportGroupUpdater) EmailDeliverySettings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.EmailDeliverySettings](
		u.rg.EmailDeliverySettings(),
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
}

// SlackDeliverySettings returns Slack delivery settings.
func (u reportGroupUpdater) SlackDeliverySettings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.SlackDeliverySettings](
		u.rg.SlackDeliverySettings(),
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
}

// TeamsDeliverySettings returns Microsoft Teams delivery settings.
func (u reportGroupUpdater) TeamsDeliverySettings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.TeamsDeliverySettings](
		u.rg.TeamsDeliverySettings(),
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
}

// ServiceNowDeliverySettings returns ServiceNow delivery settings.
func (u reportGroupUpdater) ServiceNowDeliverySettings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.ServiceNowDeliverySettings](
		u.rg.ServiceNowDeliverySettings(),
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
}

// JiraDeliverySettings returns Jira delivery settings.
func (u reportGroupUpdater) JiraDeliverySettings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.JiraDeliverySettings](
		u.rg.JiraDeliverySettings(),
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
}

// SymphonyDeliverySettings returns Symphony delivery settings.
func (u reportGroupUpdater) SymphonyDeliverySettings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.SymphonyDeliverySettings](
		u.rg.SymphonyDeliverySettings(),
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
}
