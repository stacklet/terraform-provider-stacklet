// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// RoleAssignmentResource is the model for role assignment resources.
type RoleAssignmentResource struct {
	ID        types.String `tfsdk:"id"`
	RoleName  types.String `tfsdk:"role_name"`
	Principal types.Object `tfsdk:"principal"`
	Target    types.Object `tfsdk:"target"`
}

// RoleAssignmentPrincipal is the model for a role assignment principal.
type RoleAssignmentPrincipal struct {
	Type types.String `tfsdk:"type"`
	ID   types.Int64  `tfsdk:"id"`
}

func (p RoleAssignmentPrincipal) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
		"id":   types.Int64Type,
	}
}

// RoleAssignmentTarget is the model for a role assignment target.
type RoleAssignmentTarget struct {
	Type types.String `tfsdk:"type"`
	UUID types.String `tfsdk:"uuid"`
}

func (t RoleAssignmentTarget) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
		"uuid": types.StringType,
	}
}

// Update updates the model from an API role assignment.
func (m *RoleAssignmentResource) Update(ctx context.Context, assignment *api.RoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(assignment.ID)
	m.RoleName = types.StringValue(assignment.Role.Name)

	// Update principal
	principal := assignment.GetPrincipal()
	principalAttrs := map[string]attr.Value{
		"type": types.StringValue(principal.Type),
		"id":   types.Int64Value(principal.ID),
	}
	principalObj, d := types.ObjectValue(RoleAssignmentPrincipal{}.AttributeTypes(), principalAttrs)
	diags.Append(d...)
	m.Principal = principalObj

	// Update target
	target := assignment.GetTarget()
	targetAttrs := map[string]attr.Value{
		"type": types.StringValue(target.Type),
		"uuid": types.StringPointerValue(target.UUID),
	}
	targetObj, d := types.ObjectValue(RoleAssignmentTarget{}.AttributeTypes(), targetAttrs)
	diags.Append(d...)
	m.Target = targetObj

	return diags
}

// ToAPIParams extracts the parameters needed for API calls.
func (m *RoleAssignmentResource) ToAPIParams(ctx context.Context) (string, api.RoleAssignmentPrincipal, api.RoleAssignmentTarget, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract principal
	var principal RoleAssignmentPrincipal
	d := m.Principal.As(ctx, &principal, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	// Extract target
	var target RoleAssignmentTarget
	d = m.Target.As(ctx, &target, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return "", api.RoleAssignmentPrincipal{}, api.RoleAssignmentTarget{}, diags
	}

	roleName := m.RoleName.ValueString()
	apiPrincipal := api.RoleAssignmentPrincipal{
		Type: principal.Type.ValueString(),
		ID:   principal.ID.ValueInt64(),
	}
	apiTarget := api.RoleAssignmentTarget{
		Type: target.Type.ValueString(),
		UUID: target.UUID.ValueStringPointer(),
	}

	return roleName, apiPrincipal, apiTarget, diags
}

// RoleAssignmentsDataSource is the model for role_assignments data source.
type RoleAssignmentsDataSource struct {
	Target      types.Object `tfsdk:"target"`
	Assignments types.List   `tfsdk:"assignments"`
}

// RoleAssignmentItem is a single assignment item in the data source list.
type RoleAssignmentItem struct {
	ID        types.String `tfsdk:"id"`
	RoleName  types.String `tfsdk:"role_name"`
	Principal types.Object `tfsdk:"principal"`
	Target    types.Object `tfsdk:"target"`
}

func (i RoleAssignmentItem) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"role_name": types.StringType,
		"principal": types.ObjectType{AttrTypes: RoleAssignmentPrincipal{}.AttributeTypes()},
		"target":    types.ObjectType{AttrTypes: RoleAssignmentTarget{}.AttributeTypes()},
	}
}

// Update updates the data source model from API role assignments.
func (m *RoleAssignmentsDataSource) Update(ctx context.Context, assignments []api.RoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	// Convert each assignment to a list item
	items := make([]RoleAssignmentItem, 0, len(assignments))
	for _, assignment := range assignments {
		// Build principal object
		principal := assignment.GetPrincipal()
		principalAttrs := map[string]attr.Value{
			"type": types.StringValue(principal.Type),
			"id":   types.Int64Value(principal.ID),
		}
		principalObj, d := types.ObjectValue(RoleAssignmentPrincipal{}.AttributeTypes(), principalAttrs)
		diags.Append(d...)

		// Build target object
		target := assignment.GetTarget()
		targetAttrs := map[string]attr.Value{
			"type": types.StringValue(target.Type),
			"uuid": types.StringPointerValue(target.UUID),
		}
		targetObj, d := types.ObjectValue(RoleAssignmentTarget{}.AttributeTypes(), targetAttrs)
		diags.Append(d...)

		item := RoleAssignmentItem{
			ID:        types.StringValue(assignment.ID),
			RoleName:  types.StringValue(assignment.Role.Name),
			Principal: principalObj,
			Target:    targetObj,
		}
		items = append(items, item)
	}

	// Convert to list
	assignmentsList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: RoleAssignmentItem{}.AttributeTypes()}, items)
	diags.Append(d...)
	m.Assignments = assignmentsList

	return diags
}
