// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// ConfigurationProfileServiceNowDataSource is the model for ServiceNow configuration profile data sources.
type ConfigurationProfileServiceNowDataSource struct {
	ID          types.String `tfsdk:"id"`
	Profile     types.String `tfsdk:"profile"`
	Endpoint    types.String `tfsdk:"endpoint"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	IssueType   types.String `tfsdk:"issue_type"`
	ClosedState types.String `tfsdk:"closed_state"`
}

func (m *ConfigurationProfileServiceNowDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	config := cp.Record.ServiceNowConfiguration

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.Endpoint = types.StringValue(config.Endpoint)
	m.Username = types.StringValue(config.User)
	m.Password = types.StringValue(config.Password)
	m.IssueType = types.StringValue(config.IssueType)
	m.ClosedState = types.StringValue(config.ClosedState)

	return diags
}

// ConfigurationProfileServiceNowResource is the model for ServiceNow configuration profile resources.
type ConfigurationProfileServiceNowResource struct {
	ConfigurationProfileServiceNowDataSource

	PasswordWO        types.String `tfsdk:"password_wo"`
	PasswordWOVersion types.String `tfsdk:"password_wo_version"`
}
