// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// ConfigurationProfileResourceOwnerDataSource is the model for resource owner configuration profile data sources.
type ConfigurationProfileResourceOwnerDataSource struct {
	ID           types.String `tfsdk:"id"`
	Profile      types.String `tfsdk:"profile"`
	Default      types.List   `tfsdk:"default"`
	OrgDomain    types.String `tfsdk:"org_domain"`
	OrgDomainTag types.String `tfsdk:"org_domain_tag"`
	Tags         types.List   `tfsdk:"tags"`
}

func (m *ConfigurationProfileResourceOwnerDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	config := cp.Record.ResourceOwnerConfiguration

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.Default = typehelpers.StringsList(config.Default)
	m.OrgDomain = types.StringPointerValue(config.OrgDomain)
	m.OrgDomainTag = types.StringPointerValue(config.OrgDomainTag)
	m.Tags = typehelpers.StringsList(config.Tags)

	return diags
}

// ConfigurationProfileResourceOwnerResource is the model for resource owner configuration profile resources.
type ConfigurationProfileResourceOwnerResource struct {
	ConfigurationProfileResourceOwnerDataSource
}
