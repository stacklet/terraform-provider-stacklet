// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// ConfigurationProfileSymphonyDataSource is the model for Symphony configuration profile data sources.
type ConfigurationProfileSymphonyDataSource struct {
	ID             types.String `tfsdk:"id"`
	Profile        types.String `tfsdk:"profile"`
	AgentDomain    types.String `tfsdk:"agent_domain"`
	ServiceAccount types.String `tfsdk:"service_account"`
	PrivateKey     types.String `tfsdk:"private_key"`
}

func (m *ConfigurationProfileSymphonyDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	config := cp.Record.SymphonyConfiguration

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.AgentDomain = types.StringValue(config.AgentDomain)
	m.ServiceAccount = types.StringValue(config.ServiceAccount)
	m.PrivateKey = types.StringValue(config.PrivateKey)

	return diags
}

// ConfigurationProfileSymphonyResource is the model for Symphony configuration profile resources.
type ConfigurationProfileSymphonyResource struct {
	ConfigurationProfileSymphonyDataSource

	PrivateKeyWO        types.String `tfsdk:"private_key_wo"`
	PrivateKeyWOVersion types.String `tfsdk:"private_key_wo_version"`
}
