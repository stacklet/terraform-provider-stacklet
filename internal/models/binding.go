// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// BindingResource is the model for a binding resource.
type BindingResource struct {
	BindingDataSource

	SecurityContextWO        types.String `tfsdk:"security_context_wo"`
	SecurityContextWOVersion types.String `tfsdk:"security_context_wo_version"`
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
