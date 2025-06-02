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

// BindingExecutionConfig is the model for the execution config for a binding.
type BindingExecutionConfig struct {
	Variables types.String `tfsdk:"variables"`
}

func (c BindingExecutionConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"variables": types.StringType,
	}
}
