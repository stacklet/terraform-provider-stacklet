// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
type ConfigurationProfileAccountOwnersResource ConfigurationProfileAccountOwnersDataSource
