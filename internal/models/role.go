// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

// RoleDataSource is the model for role data sources.
type RoleDataSource struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Permissions types.List   `tfsdk:"permissions"`
	System      types.Bool   `tfsdk:"system"`
}

func (m *RoleDataSource) Update(ctx context.Context, role *api.Role) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(role.ID)
	m.Name = types.StringValue(role.Name)
	m.System = types.BoolValue(role.System)

	permissions, d := types.ListValueFrom(ctx, types.StringType, role.Permissions)
	m.Permissions = permissions
	errors.AddAttributeDiags(&diags, d, "permissions")

	return diags
}
