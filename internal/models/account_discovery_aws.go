// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
