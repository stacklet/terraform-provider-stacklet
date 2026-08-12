// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// UserGroupDataSource is the model for user group data sources.
type UserGroupDataSource struct {
	ID                      types.String `tfsdk:"id"`
	UUID                    types.String `tfsdk:"uuid"`
	Name                    types.String `tfsdk:"name"`
	DisplayName             types.String `tfsdk:"display_name"`
	RoleAssignmentPrincipal types.String `tfsdk:"role_assignment_principal"`
	RoleAssignmentTarget    types.String `tfsdk:"role_assignment_target"`
}

func (m *UserGroupDataSource) Update(userGroup *api.UserGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = typehelpers.GraphQLIDValue(userGroup.ID)
	m.UUID = types.StringValue(userGroup.UUID)
	m.Name = types.StringValue(userGroup.Name)
	m.DisplayName = types.StringPointerValue(userGroup.DisplayName)
	m.RoleAssignmentPrincipal = types.StringValue(userGroup.RoleAssignmentPrincipal)
	m.RoleAssignmentTarget = types.StringValue(userGroup.RoleAssignmentTarget)

	return diags
}

// UserGroupResource is the model for user group resources.
type UserGroupResource struct {
	UserGroupDataSource
}
