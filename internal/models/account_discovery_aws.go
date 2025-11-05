// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// AccountDiscoveryAWSResource is the model for AWS account discovery resources.
type AccountDiscoveryAWSResource struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Suspended     types.Bool   `tfsdk:"suspended"`
	OrgID         types.String `tfsdk:"org_id"`
	OrgReadRole   types.String `tfsdk:"org_read_role"`
	MemberRole    types.String `tfsdk:"member_role"`
	CustodianRole types.String `tfsdk:"custodian_role"`
}

func (m *AccountDiscoveryAWSResource) Update(accountDiscovery *api.AccountDiscovery) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(accountDiscovery.ID)
	m.Name = types.StringValue(accountDiscovery.Name)
	m.Description = types.StringPointerValue(accountDiscovery.Description)
	m.Suspended = types.BoolValue(accountDiscovery.Schedule.Suspended)
	m.OrgID = types.StringValue(accountDiscovery.Config.AWSConfig.OrgID)
	m.OrgReadRole = types.StringValue(accountDiscovery.Config.AWSConfig.OrgRole)
	m.CustodianRole = types.StringValue(accountDiscovery.Config.AWSConfig.CustodianRole)

	return diags
}
