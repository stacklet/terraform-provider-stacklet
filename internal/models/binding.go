// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BindingResource is the model for a binding resource.
type BindingResource struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	Schedule             types.String `tfsdk:"schedule"`
	AccountGroupUUID     types.String `tfsdk:"account_group_uuid"`
	PolicyCollectionUUID types.String `tfsdk:"policy_collection_uuid"`
	System               types.Bool   `tfsdk:"system"`
	ExecutionConfig      types.Object `tfsdk:"execution_config"`
}

// BindingDataSource is the model for a binding data source.
type BindingDataSource BindingResource

// BindingReourceExecutionConfig is the model for the execution config for a binding data source.
type BindingResourceExecutionConfig struct {
	DryRun                   types.Bool   `tfsdk:"dry_run"`
	ResourceLimits           types.Object `tfsdk:"resource_limits"`
	SecurityContext          types.String `tfsdk:"security_context"`
	SecurityContextWO        types.String `tfsdk:"security_context_wo"`
	SecurityContextWOVersion types.String `tfsdk:"security_context_wo_version"`
	Variables                types.String `tfsdk:"variables"`
}

func (c BindingResourceExecutionConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dry_run":                     types.BoolType,
		"resource_limits":             types.ObjectType{AttrTypes: BindingExecutionConfigResourceLimits{}.AttributeTypes()},
		"security_context":            types.StringType,
		"security_context_wo":         types.StringType,
		"security_context_wo_version": types.StringType,
		"variables":                   types.StringType,
	}
}

// BindingDataSourceExecutionConfig is the model for the execution config for a binding data source.
type BindingDataSourceExecutionConfig struct {
	DryRun          types.Bool   `tfsdk:"dry_run"`
	ResourceLimits  types.Object `tfsdk:"resource_limits"`
	SecurityContext types.String `tfsdk:"security_context"`
	Variables       types.String `tfsdk:"variables"`
}

func (c BindingDataSourceExecutionConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dry_run":          types.BoolType,
		"resource_limits":  types.ObjectType{AttrTypes: BindingExecutionConfigResourceLimits{}.AttributeTypes()},
		"security_context": types.StringType,
		"variables":        types.StringType,
	}
}

// BindingExecutionConfigResourceLimits is the model for execution config resource limits for a binding.
type BindingExecutionConfigResourceLimits struct {
	Default types.Object `tfsdk:"default"`
}

func (c BindingExecutionConfigResourceLimits) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"default": types.ObjectType{AttrTypes: BindingExecutionConfigResourceLimit{}.AttributeTypes()},
	}
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
