// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileTeamsDataSource is the model for Microsoft Teams configuration profile data sources.
type ConfigurationProfileTeamsDataSource struct {
	ID       types.String `tfsdk:"id"`
	Profile  types.String `tfsdk:"profile"`
	Webhooks types.List   `tfsdk:"webhook"`
}

// ConfigurationProfileTeamsResource is the model for Microsoft Teams configuration profile resources.
type ConfigurationProfileTeamsResource ConfigurationProfileTeamsDataSource

// TeamsWebhook is the model for a Microsoft Teams webhook.
type TeamsWebhook struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func (w TeamsWebhook) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"url":  types.StringType,
	}
}

// TeamsWebhookWithSecret is the model for a Microsoft Teams webhook including the URL as secret.
type TeamsWebhookWithSecret struct {
	TeamsWebhook

	URLWO        types.String `tfsdk:"url_wo"`
	URLWOVersion types.String `tfsdk:"url_wo_version"`
}

func (w TeamsWebhookWithSecret) AttributeTypes() map[string]attr.Type {
	attrs := w.TeamsWebhook.AttributeTypes()
	attrs["url_wo"] = types.StringType
	attrs["url_wo_version"] = types.StringType
	return attrs
}
