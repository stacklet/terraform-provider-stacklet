// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
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
	m.UserFields = typehelpers.StringsList(slackConfig.UserFields)

	webhooks, d := typehelpers.ObjectList[SlackWebhook](
		cp.Record.SlackConfiguration.Webhooks,
		func(entry api.SlackWebhook) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name": types.StringValue(entry.Name),
				"url":  types.StringValue(entry.URL),
			}, nil
		},
	)
	m.Webhooks = webhooks
	errors.AddAttributeDiags(&diags, d, "webhook")

	return diags
}

// ConfigurationProfileSlackResource is the model for Slack configuration profile resources.
type ConfigurationProfileSlackResource struct {
	ConfigurationProfileSlackDataSource

	TokenWO        types.String `tfsdk:"token_wo"`
	TokenWOVersion types.String `tfsdk:"token_wo_version"`
}

func (m *ConfigurationProfileSlackResource) Update(ctx context.Context, cp api.ConfigurationProfile, webhookVersions map[string]string) diag.Diagnostics {
	var diags diag.Diagnostics

	// fetch current webhook names to preserve declared order
	webhookNames, d := typehelpers.ListItemsIdentifiers(m.Webhooks, "name")
	errors.AddAttributeDiags(&diags, d, "webhook")

	d = m.ConfigurationProfileSlackDataSource.Update(cp)
	diags.Append(d...)

	if !m.Webhooks.IsNull() {
		// sort entries according to keep previous ordering
		if webhookNames != nil {
			webhooks, d := typehelpers.ListSortedEntries[SlackWebhook](m.Webhooks, "name", webhookNames)
			m.Webhooks = webhooks
			errors.AddAttributeDiags(&diags, d, "webhook")
		}

		// extend entries with secrets
		elems := m.Webhooks.Elements()
		values := make([]attr.Value, 0, len(elems))
		for _, entry := range elems {
			obj, _ := entry.(types.Object)

			name, d := typehelpers.ObjectStringIdentifier(obj, "name")
			errors.AddAttributeDiags(&diags, d, "webhook")
			if diags.HasError() {
				return diags
			}

			var woVersion types.String
			if version, ok := webhookVersions[name]; ok {
				woVersion = types.StringValue(version)
			} else {
				woVersion = types.StringNull()
			}

			webhook, d := typehelpers.UpdatedObject(
				ctx,
				obj,
				map[string]attr.Value{
					"url_wo":         types.StringNull(), // always empty since it's not stored in the state
					"url_wo_version": woVersion,
				},
			)
			values = append(values, webhook)
			errors.AddAttributeDiags(&diags, d, "webhook")
		}

		webhooks, d := types.ListValue(types.ObjectType{AttrTypes: SlackWebhookWithSecret{}.AttributeTypes()}, values)
		m.Webhooks = webhooks
		errors.AddAttributeDiags(&diags, d, "webhook")
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
