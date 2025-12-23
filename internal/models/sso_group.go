// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// SSOGroupDataSource is the model for SSO group data sources.
type SSOGroupDataSource struct {
	ID                      types.String `tfsdk:"id"`
	DisplayName             types.String `tfsdk:"display_name"`
	Name                    types.String `tfsdk:"name"`
	RoleAssignmentPrincipal types.String `tfsdk:"role_assignment_principal"`
}

func (m *SSOGroupDataSource) Update(ctx context.Context, ssoGroup *api.SSOGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(ssoGroup.ID)
	m.DisplayName = types.StringPointerValue(ssoGroup.DisplayName)
	m.Name = types.StringValue(ssoGroup.Name)
	m.RoleAssignmentPrincipal = types.StringValue(ssoGroup.RoleAssignmentPrincipal)

	return diags
}
