// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

// ConfigurationProfileSlackDataSource is the model for Slack configuration profile data sources.
type ConfigurationProfileSlackDataSource struct {
	ID         types.String `tfsdk:"id"`
	Profile    types.String `tfsdk:"profile"`
	Token      types.String `tfsdk:"token"`
	UserFields types.List   `tfsdk:"user_fields"`
	Webhooks   types.List   `tfsdk:"webhook"`
}

func (m *ConfigurationProfileSlackDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	slackConfig := cp.Record.SlackConfiguration

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.Token = types.StringPointerValue(slackConfig.Token)
	m.UserFields = tftypes.StringsList(slackConfig.UserFields)

	webhooks, d := tftypes.ObjectList[SlackWebhook](
		cp.Record.SlackConfiguration.Webhooks,
		func(entry api.SlackWebhook) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name": types.StringValue(entry.Name),
				"url":  types.StringValue(entry.URL),
			}, nil
		},
	)
	m.Webhooks = webhooks
	diags.Append(d...)

	return diags
}

// ConfigurationProfileSlackResource is the model for Slack configuration profile resources.
type ConfigurationProfileSlackResource struct {
	ConfigurationProfileSlackDataSource

	TokenWO        types.String `tfsdk:"token_wo"`
	TokenWOVersion types.String `tfsdk:"token_wo_version"`
}

func (m *ConfigurationProfileSlackResource) Update(ctx context.Context, cp api.ConfigurationProfile, webhookVersions map[string]string) diag.Diagnostics {
	// fetch current webhook names to preserve declared order
	webhookNames := tftypes.ListItemsIdentifiers(m.Webhooks, "name")

	diags := m.ConfigurationProfileSlackDataSource.Update(cp)

	if !m.Webhooks.IsNull() {

		// sort entries according to keep previous ordering
		if webhookNames != nil {
			webhooks, d := tftypes.ListSortedEntries[SlackWebhook](m.Webhooks, "name", webhookNames)
			m.Webhooks = webhooks
			diags.Append(d...)
		}

		// extend entries with secrets
		elems := m.Webhooks.Elements()
		values := make([]attr.Value, 0, len(elems))
		for _, entry := range elems {
			obj, _ := entry.(types.Object)

			var woVersion types.String
			if version, ok := webhookVersions[tftypes.ObjectStringIdentifier(obj, "name")]; ok {
				woVersion = types.StringValue(version)
			} else {
				woVersion = types.StringNull()
			}

			webhook, d := tftypes.UpdatedObject(
				ctx,
				obj,
				map[string]attr.Value{
					"url_wo":         types.StringNull(), // always empty since it's not stored in the state
					"url_wo_version": woVersion,
				},
			)
			values = append(values, webhook)
			diags.Append(d...)
		}

		webhooks, d := types.ListValue(types.ObjectType{AttrTypes: SlackWebhookWithSecret{}.AttributeTypes()}, values)
		m.Webhooks = webhooks
		diags.Append(d...)
	}

	return diags
}

// SlackWebhook is the model for a Slack webhook.
type SlackWebhook struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func (w SlackWebhook) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"url":  types.StringType,
	}
}

// SlackWebhookWithSecret is the model for a Slack webhook including the URL as secret.
type SlackWebhookWithSecret struct {
	SlackWebhook

	URLWO        types.String `tfsdk:"url_wo"`
	URLWOVersion types.String `tfsdk:"url_wo_version"`
}

func (w SlackWebhookWithSecret) AttributeTypes() map[string]attr.Type {
	attrs := w.SlackWebhook.AttributeTypes()
	attrs["url_wo"] = types.StringType
	attrs["url_wo_version"] = types.StringType
	return attrs
}
