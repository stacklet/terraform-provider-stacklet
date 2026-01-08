// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// UserDataSource is the model for user data sources.
type UserDataSource struct {
	ID                      types.String `tfsdk:"id"`
	Active                  types.Bool   `tfsdk:"active"`
	DisplayName             types.String `tfsdk:"display_name"`
	Email                   types.String `tfsdk:"email"`
	Name                    types.String `tfsdk:"name"`
	RoleAssignmentPrincipal types.String `tfsdk:"role_assignment_principal"`
	SSOUser                 types.Bool   `tfsdk:"sso_user"`
	Username                types.String `tfsdk:"username"`
}

func (m *UserDataSource) Update(ctx context.Context, user *api.User) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(user.ID)
	m.Active = types.BoolValue(user.Active)
	m.DisplayName = types.StringPointerValue(user.DisplayName)
	m.Email = types.StringPointerValue(user.Email)
	m.Name = types.StringPointerValue(user.Name)
	m.RoleAssignmentPrincipal = types.StringValue(user.RoleAssignmentPrincipal)
	m.SSOUser = types.BoolValue(user.SSOUser)
	m.Username = types.StringPointerValue(user.Username)

	return diags
}
