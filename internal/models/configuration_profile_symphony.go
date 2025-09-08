// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileSymphonyDataSource is the model for Symphony configuration profile data sources.
type ConfigurationProfileSymphonyDataSource struct {
	ID             types.String `tfsdk:"id"`
	Profile        types.String `tfsdk:"profile"`
	AgentDomain    types.String `tfsdk:"agent_domain"`
	ServiceAccount types.String `tfsdk:"service_account"`
	PrivateKey     types.String `tfsdk:"private_key"`
}

// ConfigurationProfileSymphonyResource is the model for Symphony configuration profile resources.
type ConfigurationProfileSymphonyResource struct {
	ConfigurationProfileSymphonyDataSource

	PrivateKeyWO        types.String `tfsdk:"private_key_wo"`
	PrivateKeyWOVersion types.String `tfsdk:"private_key_wo_version"`
}
