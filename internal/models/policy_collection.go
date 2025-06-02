// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyCollectionResource is the model for a policy collection resource.
type PolicyCollectionResource struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	AutoUpdate    types.Bool   `tfsdk:"auto_update"`
	System        types.Bool   `tfsdk:"system"`
	Dynamic       types.Bool   `tfsdk:"dynamic"`
	DynamicConfig types.Object `tfsdk:"dynamic_config"`
}

// PolicyCollectionDatasource is the model for a policy collection data source.
type PolicyCollectionDataSource PolicyCollectionResource

// PolicyCollectionDynamicConfig is the model for the dynamic configuration for a policy collection.
type PolicyCollectionDynamicConfig struct {
	RepositoryUUID     types.String `tfsdk:"repository_uuid"`
	Namespace          types.String `tfsdk:"namespace"`
	BranchName         types.String `tfsdk:"branch_name"`
	PolicyDirectories  types.List   `tfsdk:"policy_directories"`
	PolicyFileSuffixes types.List   `tfsdk:"policy_file_suffixes"`
}

func (c PolicyCollectionDynamicConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository_uuid":      types.StringType,
		"namespace":            types.StringType,
		"branch_name":          types.StringType,
		"policy_directories":   types.ListType{ElemType: types.StringType},
		"policy_file_suffixes": types.ListType{ElemType: types.StringType},
	}
}
