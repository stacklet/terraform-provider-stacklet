// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// BindingDataSource is the model for a binding data source.
type BindingDataSource struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	Schedule             types.String `tfsdk:"schedule"`
	AccountGroupUUID     types.String `tfsdk:"account_group_uuid"`
	PolicyCollectionUUID types.String `tfsdk:"policy_collection_uuid"`
	System               types.Bool   `tfsdk:"system"`
	DryRun               types.Bool   `tfsdk:"dry_run"`
	ResourceLimits       types.Object `tfsdk:"resource_limits"`
	PolicyResourceLimits types.List   `tfsdk:"policy_resource_limit"`
	SecurityContext      types.String `tfsdk:"security_context"`
	Variables            types.String `tfsdk:"variables"`
}

func (m *BindingDataSource) Update(ctx context.Context, binding *api.Binding) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(binding.ID)
	m.UUID = types.StringValue(binding.UUID)
	m.Name = types.StringValue(binding.Name)
	m.Description = types.StringPointerValue(binding.Description)
	m.AutoDeploy = types.BoolValue(binding.AutoDeploy)
	m.Schedule = types.StringPointerValue(binding.Schedule)
	m.AccountGroupUUID = types.StringValue(binding.AccountGroup.UUID)
	m.PolicyCollectionUUID = types.StringValue(binding.PolicyCollection.UUID)
	m.System = types.BoolValue(binding.System)
	m.DryRun = types.BoolPointerValue(binding.DryRun())
	m.SecurityContext = types.StringPointerValue(binding.SecurityContext())

	variablesString, d := typehelpers.JSONString(binding.ExecutionConfig.Variables)
	errors.AddAttributeDiags(&diags, d, "variables")
	m.Variables = variablesString

	defLimit := binding.DefaultResourceLimits()
	defaultLimits, d := typehelpers.ObjectValue(
		ctx,
		defLimit,
		func() (*BindingExecutionConfigResourceLimit, diag.Diagnostics) {
			return &BindingExecutionConfigResourceLimit{
				MaxCount:      types.Int32PointerValue(defLimit.MaxCount),
				MaxPercentage: types.Float32PointerValue(defLimit.MaxPercentage),
				RequiresBoth:  types.BoolValue(defLimit.RequiresBoth),
			}, nil
		},
	)
	errors.AddAttributeDiags(&diags, d, "resource_limits")
	m.ResourceLimits = defaultLimits

	policyLimits, d := typehelpers.ObjectList[BindingExecutionConfigPolicyResourceLimit](
		binding.PolicyResourceLimits(),
		func(entry api.BindingExecutionConfigResourceLimitsPolicyOverrides) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"policy_name":    types.StringValue(entry.PolicyName),
				"max_count":      types.Int32PointerValue(entry.Limit.MaxCount),
				"max_percentage": types.Float32PointerValue(entry.Limit.MaxPercentage),
				"requires_both":  types.BoolValue(entry.Limit.RequiresBoth),
			}, nil
		},
	)
	errors.AddAttributeDiags(&diags, d, "policy_resource_limit")
	m.PolicyResourceLimits = policyLimits

	return diags
}

// BindingResource is the model for a binding resource.
type BindingResource struct {
	BindingDataSource

	SecurityContextWO        types.String `tfsdk:"security_context_wo"`
	SecurityContextWOVersion types.String `tfsdk:"security_context_wo_version"`
}

func (m *BindingResource) Update(ctx context.Context, binding *api.Binding) diag.Diagnostics {
	// Save original values for resource-specific logic
	originalVariables := m.Variables
	originalResourceLimits := m.ResourceLimits

	diags := m.BindingDataSource.Update(ctx, binding)

	// the API returns an empty dictionary for both null or empty strings. In
	// that case don't modify the expected value.
	if m.Variables.ValueString() == "{}" {
		m.Variables = originalVariables
	}

	// if default resource limits are set in the config but empty, apply default
	if binding.DefaultResourceLimits() == nil && !originalResourceLimits.IsNull() {
		var d diag.Diagnostics
		defLimits := BindingExecutionConfigResourceLimit{RequiresBoth: types.BoolValue(false)}
		resourceLimits, d := types.ObjectValueFrom(ctx, defLimits.AttributeTypes(), &defLimits)
		errors.AddAttributeDiags(&diags, d, "resource_limits")
		m.ResourceLimits = resourceLimits
	}

	return diags
}

// BindingExecutionConfigResourceLimit is the model for a resource limit in binding execution config.
type BindingExecutionConfigResourceLimit struct {
	MaxCount      types.Int32   `tfsdk:"max_count"`
	MaxPercentage types.Float32 `tfsdk:"max_percentage"`
	RequiresBoth  types.Bool    `tfsdk:"requires_both"`
}

func (c BindingExecutionConfigResourceLimit) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"max_count":      types.Int32Type,
		"max_percentage": types.Float32Type,
		"requires_both":  types.BoolType,
	}
}

// BindingExecutionConfigPolicyResourceLimit is the model for a policy resource limit in binding execution config.
type BindingExecutionConfigPolicyResourceLimit struct {
	BindingExecutionConfigResourceLimit

	PolicyName types.String `tfsdk:"policy_name"`
}

func (c BindingExecutionConfigPolicyResourceLimit) AttributeTypes() map[string]attr.Type {
	m := c.BindingExecutionConfigResourceLimit.AttributeTypes()
	m["policy_name"] = types.StringType
	return m
}
