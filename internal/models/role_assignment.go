// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// RoleAssignmentResource is the model for role assignment resources.
// Principal and Target are opaque string identifiers.
type RoleAssignmentResource struct {
	ID        types.String `tfsdk:"id"`
	RoleName  types.String `tfsdk:"role_name"`
	Principal types.String `tfsdk:"principal"`
	Target    types.String `tfsdk:"target"`
}

// Update updates the model from an API role assignment.
func (m *RoleAssignmentResource) Update(ctx context.Context, assignment *api.RoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(assignment.ID)
	m.RoleName = types.StringValue(assignment.Role.Name)
	m.Principal = types.StringValue(assignment.GetPrincipal())
	m.Target = types.StringValue(assignment.GetTarget())

	return diags
}

// ToAPIParams extracts the parameters needed for API calls.
// Returns roleName, principal (opaque string), target (opaque string).
func (m *RoleAssignmentResource) ToAPIParams(ctx context.Context) (string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics

	roleName := m.RoleName.ValueString()
	principal := m.Principal.ValueString()
	target := m.Target.ValueString()

	return roleName, principal, target, diags
}

// RoleAssignmentsDataSource is the model for role_assignments data source.
type RoleAssignmentsDataSource struct {
	Target      types.String `tfsdk:"target"`
	Assignments types.List   `tfsdk:"assignments"`
}

// RoleAssignmentItem is a single assignment item in the data source list.
type RoleAssignmentItem struct {
	ID        types.String `tfsdk:"id"`
	RoleName  types.String `tfsdk:"role_name"`
	Principal types.String `tfsdk:"principal"`
	Target    types.String `tfsdk:"target"`
}

func (i RoleAssignmentItem) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"role_name": types.StringType,
		"principal": types.StringType,
		"target":    types.StringType,
	}
}

// Update updates the data source model from API role assignments.
func (m *RoleAssignmentsDataSource) Update(ctx context.Context, assignments []api.RoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	// Convert each assignment to a list item
	items := make([]RoleAssignmentItem, 0, len(assignments))
	for _, assignment := range assignments {
		item := RoleAssignmentItem{
			ID:        types.StringValue(assignment.ID),
			RoleName:  types.StringValue(assignment.Role.Name),
			Principal: types.StringValue(assignment.GetPrincipal()),
			Target:    types.StringValue(assignment.GetTarget()),
		}
		items = append(items, item)
	}

	// Convert to list
	assignmentsList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: RoleAssignmentItem{}.AttributeTypes()}, items)
	diags.Append(d...)
	m.Assignments = assignmentsList

	return diags
}
