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
	UserFields types.List   `tfsdk:"user_fields"`
	Webhooks   types.List   `tfsdk:"webhook"`
}

// SlackWebhook is the model for a Microsoft Slack webhook.
type SlackWebhook struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func (w SlackWebhook) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}
