// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
