// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BindingResource struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	Schedule             types.String `tfsdk:"schedule"`
	Variables            types.String `tfsdk:"variables"`
	AccountGroupUUID     types.String `tfsdk:"account_group_uuid"`
	PolicyCollectionUUID types.String `tfsdk:"policy_collection_uuid"`
	System               types.Bool   `tfsdk:"system"`
}

type BindingDataSource struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	Schedule             types.String `tfsdk:"schedule"`
	Variables            types.String `tfsdk:"variables"`
	AccountGroupUUID     types.String `tfsdk:"account_group_uuid"`
	PolicyCollectionUUID types.String `tfsdk:"policy_collection_uuid"`
	System               types.Bool   `tfsdk:"system"`
}
