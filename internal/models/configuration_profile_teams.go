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

// TeamsWebhook is the model for a Microsoft Teams webhook.
type TeamsWebhook struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func (w TeamsWebhook) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}
