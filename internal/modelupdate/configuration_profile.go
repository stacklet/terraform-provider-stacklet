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

// NewConfigurationProfileUpdater returns a configuration profile updater helper.
func NewConfigurationProfileUpdater(cp api.ConfigurationProfile) configurationProfileUpdater {
	return configurationProfileUpdater{cp: cp}
}

type configurationProfileUpdater struct {
	cp api.ConfigurationProfile
}

// JiraProjects returns a list of Jira projects settings.
func (u configurationProfileUpdater) JiraProjects() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.JiraProject](
		u.cp.Record.JiraConfiguration.Projects,
		func(entry api.JiraProject) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"closed_status": types.StringValue(entry.ClosedStatus),
				"issue_type":    types.StringValue(entry.IssueType),
				"name":          types.StringValue(entry.Name),
				"project":       types.StringValue(entry.Project),
			}, nil
		},
	)
}

// TeamsWebhooks returns a list of Microsoft Teams webhooks.
func (u configurationProfileUpdater) TeamsWebhooks() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.TeamsWebhook](
		u.cp.Record.TeamsConfiguration.Webhooks,
		func(entry api.TeamsWebhook) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name": types.StringValue(entry.Name),
			}, nil
		},
	)
}

// SlackWebhooks returns a list of Slack webhooks.
func (u configurationProfileUpdater) SlackWebhooks() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.SlackWebhook](
		u.cp.Record.SlackConfiguration.Webhooks,
		func(entry api.SlackWebhook) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name": types.StringValue(entry.Name),
			}, nil
		},
	)
}

// AccountOwnersDefault returns a list of account owner defaults.
func (u configurationProfileUpdater) AccountOwnersDefault() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.AccountOwners](
		u.cp.Record.AccountOwnersConfiguration.Default,
		func(entry api.AccountOwners) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"account": types.StringValue(entry.Account),
				"owners":  tftypes.StringsList(entry.Owners),
			}, nil
		},
	)
}
