// Copyright (c) 2025 - Stacklet, Inc.

package modelupdate

import (
	"fmt"

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

// JiraProjects returns a list of Jira projects settings, optionally in the order specified by names.
func (u configurationProfileUpdater) JiraProjects(names []string) (basetypes.ListValue, diag.Diagnostics) {
	projectValueAttrs := func(proj api.JiraProject) map[string]attr.Value {
		return map[string]attr.Value{
			"closed_status": types.StringValue(proj.ClosedStatus),
			"issue_type":    types.StringValue(proj.IssueType),
			"name":          types.StringValue(proj.Name),
			"project":       types.StringValue(proj.Project),
		}
	}

	if names == nil {
		return tftypes.ObjectList[models.JiraProject](
			u.cp.Record.JiraConfiguration.Projects,
			func(entry api.JiraProject) (map[string]attr.Value, diag.Diagnostics) {
				return projectValueAttrs(entry), nil
			},
		)
	}

	var diags diag.Diagnostics

	projects := map[string]api.JiraProject{}
	for _, proj := range u.cp.Record.JiraConfiguration.Projects {
		projects[proj.Name] = proj
	}
	attrTypes := models.JiraProject{}.AttributeTypes()

	values := []attr.Value{}
	for _, name := range names {
		proj, ok := projects[name]
		if !ok {
			diags.AddError(
				"Project entry not found",
				fmt.Sprintf("Project entry '%s' not found in API result", name),
			)
			return basetypes.ListValue{}, diags
		}

		value, diags := types.ObjectValue(attrTypes, projectValueAttrs(proj))
		if diags.HasError() {
			return basetypes.ListValue{}, diags
		}
		values = append(values, value)
	}

	return types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
}

// TeamsWebhooks returns a list of Microsoft Teams webhooks.
func (u configurationProfileUpdater) TeamsWebhooks() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.TeamsWebhook](
		u.cp.Record.TeamsConfiguration.Webhooks,
		func(entry api.TeamsWebhook) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name": types.StringValue(entry.Name),
				"url":  types.StringValue(entry.URL),
			}, nil
		},
	)
}

// TeamsWebhooksWithSecrets returns a list of Microsoft Teams webhooks with the secret fields, optionally in the order specified by names.
func (u configurationProfileUpdater) TeamsWebhooksWithSecret(versions map[string]string, names []string) (basetypes.ListValue, diag.Diagnostics) {
	webhookValueAttrs := func(wh api.TeamsWebhook) map[string]attr.Value {
		var woVersion basetypes.StringValue
		if version, ok := versions[wh.Name]; ok {
			woVersion = types.StringValue(version)
		} else {
			woVersion = types.StringNull()
		}

		return map[string]attr.Value{
			"name":           types.StringValue(wh.Name),
			"url":            types.StringValue(wh.URL),
			"url_wo":         types.StringNull(), // always empty since it's not stored in the state
			"url_wo_version": woVersion,
		}
	}

	if names == nil {
		return tftypes.ObjectList[models.TeamsWebhookWithSecret](
			u.cp.Record.TeamsConfiguration.Webhooks,
			func(entry api.TeamsWebhook) (map[string]attr.Value, diag.Diagnostics) {
				return webhookValueAttrs(entry), nil
			},
		)
	}

	var diags diag.Diagnostics

	webhooks := map[string]api.TeamsWebhook{}
	for _, wh := range u.cp.Record.TeamsConfiguration.Webhooks {
		webhooks[wh.Name] = wh
	}
	attrTypes := models.TeamsWebhookWithSecret{}.AttributeTypes()

	values := []attr.Value{}
	for _, name := range names {
		wh, ok := webhooks[name]
		if !ok {
			diags.AddError(
				"Webhook entry not found",
				fmt.Sprintf("Webhook entry '%s' not found in API result", name),
			)
			return basetypes.ListValue{}, diags
		}

		value, diags := types.ObjectValue(attrTypes, webhookValueAttrs(wh))
		if diags.HasError() {
			return basetypes.ListValue{}, diags
		}
		values = append(values, value)
	}

	return types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
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
