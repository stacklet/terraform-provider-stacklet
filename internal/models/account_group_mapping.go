// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// AccountGroupMappingResource is the model for an account group mapping resource.
type AccountGroupMappingResource struct {
	ID         types.String `tfsdk:"id"`
	GroupUUID  types.String `tfsdk:"group_uuid"`
	AccountKey types.String `tfsdk:"account_key"`
}

func (m *AccountGroupMappingResource) Update(accountGroupMapping *api.AccountGroupMapping) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(accountGroupMapping.ID)
	m.GroupUUID = types.StringValue(accountGroupMapping.GroupUUID)
	m.AccountKey = types.StringValue(accountGroupMapping.AccountKey)

	return diags
}
