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

// SlackWebhooks returns a list of Slack webhooks.
func (u configurationProfileUpdater) SlackWebhooks() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.SlackWebhook](
		u.cp.Record.SlackConfiguration.Webhooks,
		func(entry api.SlackWebhook) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name": types.StringValue(entry.Name),
				"url":  types.StringValue(entry.URL),
			}, nil
		},
	)
}

// SlackWebhooksWithSecret returns a list of Slack webhooks with the secret fields, optionally in the order specified by names.
func (u configurationProfileUpdater) SlackWebhooksWithSecret(versions map[string]string, names []string) (basetypes.ListValue, diag.Diagnostics) {
	webhookValueAttrs := func(wh api.SlackWebhook) map[string]attr.Value {
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
		return tftypes.ObjectList[models.SlackWebhookWithSecret](
			u.cp.Record.SlackConfiguration.Webhooks,
			func(entry api.SlackWebhook) (map[string]attr.Value, diag.Diagnostics) {
				return webhookValueAttrs(entry), nil
			},
		)
	}

	var diags diag.Diagnostics

	webhooks := map[string]api.SlackWebhook{}
	for _, wh := range u.cp.Record.SlackConfiguration.Webhooks {
		webhooks[wh.Name] = wh
	}
	attrTypes := models.SlackWebhookWithSecret{}.AttributeTypes()

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

// MSTeamsAccessConfig returns the Microsoft Teams access configuration as an object.
func (u configurationProfileUpdater) MSTeamsAccessConfig() (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	cfg := u.cp.Record.MSTeamsConfiguration.AccessConfig
	if cfg == nil {
		return types.ObjectNull(models.MSTeamsAccessConfig{}.AttributeTypes()), diags
	}

	botApplication, d := types.ObjectValue(
		models.MSTeamsBotApplication{}.AttributeTypes(),
		map[string]attr.Value{
			"download_url": types.StringValue(cfg.BotApplication.DownloadURL),
			"version":      types.StringValue(cfg.BotApplication.Version),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	var publishedApplication basetypes.ObjectValue
	if cfg.PublishedApplication == nil || (cfg.PublishedApplication.CatalogID == nil && cfg.PublishedApplication.Version == nil) {
		publishedApplication = types.ObjectNull(models.MSTeamsPublishedApplication{}.AttributeTypes())
	} else {
		var d diag.Diagnostics
		publishedApplication, d = types.ObjectValue(
			models.MSTeamsPublishedApplication{}.AttributeTypes(),
			map[string]attr.Value{
				"catalog_id": types.StringPointerValue(cfg.PublishedApplication.CatalogID),
				"version":    types.StringPointerValue(cfg.PublishedApplication.Version),
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}
	}

	return types.ObjectValue(
		models.MSTeamsAccessConfig{}.AttributeTypes(),
		map[string]attr.Value{
			"client_id":             types.StringValue(cfg.ClientID),
			"roundtrip_digest":      types.StringValue(cfg.RoundtripDigest),
			"tenant_id":             types.StringValue(cfg.TenantID),
			"bot_application":       botApplication,
			"published_application": publishedApplication,
		},
	)
}

// MSTeamsCustomerConfig returns the Microsoft Teams customer configuration as an object.
func (u configurationProfileUpdater) MSTeamsCustomerConfig() (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	cfg := u.cp.Record.MSTeamsConfiguration.CustomerConfig

	tags, d := cfg.Tags.TagsMap()
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	var version types.String
	if cfg.TerraformModule.Version != nil {
		version = types.StringValue(*cfg.TerraformModule.Version)
	} else {
		version = types.StringNull()
	}

	terraformModule, d := types.ObjectValue(
		models.TerraformModule{}.AttributeTypes(),
		map[string]attr.Value{
			"repository_url": types.StringValue(cfg.TerraformModule.RepositoryURL),
			"source":         types.StringValue(cfg.TerraformModule.Source),
			"version":        version,
			"variables_json": types.StringValue(cfg.TerraformModule.VariablesJSON),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return types.ObjectValue(
		models.MSTeamsCustomerConfig{}.AttributeTypes(),
		map[string]attr.Value{
			"bot_endpoint":     types.StringValue(cfg.BotEndpoint),
			"oidc_client":      types.StringValue(cfg.OIDCClient),
			"oidc_issuer":      types.StringValue(cfg.OIDCIssuer),
			"prefix":           types.StringValue(cfg.Prefix),
			"roundtrip_digest": types.StringValue(cfg.RoundtripDigest),
			"tags":             tags,
			"terraform_module": terraformModule,
		},
	)
}

// MSTeamsChannelMappings returns a list of Microsoft Teams channel mappings.
func (u configurationProfileUpdater) MSTeamsChannelMappings() (basetypes.ListValue, diag.Diagnostics) {
	return tftypes.ObjectList[models.MSTeamsChannelMapping](
		u.cp.Record.MSTeamsConfiguration.ChannelMappings,
		func(entry api.MSTeamsChannelMapping) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name":       types.StringValue(entry.Name),
				"team_id":    types.StringValue(string(entry.TeamID)),
				"channel_id": types.StringValue(entry.ChannelID),
			}, nil
		},
	)
}

// MSTeamsBotApplication returns the Microsoft Teams bot application as an object.
func (u configurationProfileUpdater) MSTeamsBotApplication() (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	cfg := u.cp.Record.MSTeamsConfiguration.AccessConfig

	if cfg == nil {
		return types.ObjectNull(models.MSTeamsBotApplication{}.AttributeTypes()), diags
	}

	return types.ObjectValue(
		models.MSTeamsBotApplication{}.AttributeTypes(),
		map[string]attr.Value{
			"download_url": types.StringValue(cfg.BotApplication.DownloadURL),
			"version":      types.StringValue(cfg.BotApplication.Version),
		},
	)
}

// MSTeamsPublishedApplication returns the Microsoft Teams published application as an object.
func (u configurationProfileUpdater) MSTeamsPublishedApplication() (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	cfg := u.cp.Record.MSTeamsConfiguration.AccessConfig

	if cfg == nil || cfg.PublishedApplication == nil || (cfg.PublishedApplication.CatalogID == nil && cfg.PublishedApplication.Version == nil) {
		return types.ObjectNull(models.MSTeamsPublishedApplication{}.AttributeTypes()), diags
	}

	return types.ObjectValue(
		models.MSTeamsPublishedApplication{}.AttributeTypes(),
		map[string]attr.Value{
			"catalog_id": types.StringValue(*cfg.PublishedApplication.CatalogID),
			"version":    types.StringValue(*cfg.PublishedApplication.Version),
		},
	)
}

// MSTeamsEntityDetails returns the Microsoft Teams entity details as an object.
func (u configurationProfileUpdater) MSTeamsEntityDetails() (basetypes.ObjectValue, diag.Diagnostics) {
	channels, diags := tftypes.ObjectList[models.MSTeamsChannelDetails](
		u.cp.Record.MSTeamsConfiguration.EntityDetails.Channels,
		func(entry api.MSTeamsChannelDetail) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"id":   types.StringValue(entry.ID),
				"name": types.StringValue(entry.Name),
			}, nil
		},
	)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	teams, teamsDiags := tftypes.ObjectList[models.MSTeamsTeamDetails](
		u.cp.Record.MSTeamsConfiguration.EntityDetails.Teams,
		func(entry api.MSTeamsTeamDetail) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"id":   types.StringValue(entry.ID),
				"name": types.StringValue(entry.Name),
			}, nil
		},
	)
	diags.Append(teamsDiags...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return types.ObjectValue(
		models.MSTeamsEntityDetails{}.AttributeTypes(),
		map[string]attr.Value{
			"channels": channels,
			"teams":    teams,
		},
	)
}
