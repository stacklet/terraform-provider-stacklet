// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// ConfigurationProfileAccountOwnersDataSource is the model for account owners configuration profile data sources.
type ConfigurationProfileAccountOwnersDataSource struct {
	ID           types.String `tfsdk:"id"`
	Profile      types.String `tfsdk:"profile"`
	Default      types.List   `tfsdk:"default"`
	OrgDomain    types.String `tfsdk:"org_domain"`
	OrgDomainTag types.String `tfsdk:"org_domain_tag"`
	Tags         types.List   `tfsdk:"tags"`
}

func (m *ConfigurationProfileAccountOwnersDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	config := cp.Record.AccountOwnersConfiguration

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.OrgDomain = types.StringPointerValue(config.OrgDomain)
	m.OrgDomainTag = types.StringPointerValue(config.OrgDomainTag)
	m.Tags = typehelpers.StringsList(config.Tags)

	defaultOwners, d := typehelpers.ObjectList[AccountOwners](
		config.Default,
		func(entry api.AccountOwners) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"account": types.StringValue(entry.Account),
				"owners":  typehelpers.StringsList(entry.Owners),
			}, nil
		},
	)
	errors.AddAttributeDiags(&diags, d, "default")
	m.Default = defaultOwners

	return diags
}

// AccountOwners is the model for account owners.
type AccountOwners struct {
	Account types.String `tfsdk:"account"`
	Owners  types.List   `tfsdk:"owners"`
}

func (o AccountOwners) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account": types.StringType,
		"owners":  types.ListType{ElemType: types.StringType},
	}
}

// ConfigurationProfileAccountOwnersResource is the model for account owners configuration profile resources.
type ConfigurationProfileAccountOwnersResource struct {
	ConfigurationProfileAccountOwnersDataSource
}
