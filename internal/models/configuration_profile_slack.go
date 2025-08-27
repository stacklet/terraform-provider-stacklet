// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileSlackDataSource is the model for Slack configuration profile data sources.
type ConfigurationProfileSlackDataSource struct {
	ID         types.String `tfsdk:"id"`
	Profile    types.String `tfsdk:"profile"`
	Token      types.String `tfsdk:"token"`
	UserFields types.List   `tfsdk:"user_fields"`
	Webhooks   types.List   `tfsdk:"webhook"`
}

// ConfigurationProfileSlackResource is the model for Slack configuration profile resources.
type ConfigurationProfileSlackResource struct {
	ConfigurationProfileSlackDataSource

	TokenWO        types.String `tfsdk:"token_wo"`
	TokenWOVersion types.String `tfsdk:"token_wo_version"`
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
