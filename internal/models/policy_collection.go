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

// PolicyCollectionResource is the model for a policy collection resource.
type PolicyCollectionResource struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	CloudProvider        types.String `tfsdk:"cloud_provider"`
	AutoUpdate           types.Bool   `tfsdk:"auto_update"`
	System               types.Bool   `tfsdk:"system"`
	Dynamic              types.Bool   `tfsdk:"dynamic"`
	DynamicConfig        types.Object `tfsdk:"dynamic_config"`
	RoleAssignmentTarget types.String `tfsdk:"role_assignment_target"`
}

func (m *PolicyCollectionResource) Update(ctx context.Context, policyCollection *api.PolicyCollection) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(policyCollection.ID)
	m.UUID = types.StringValue(policyCollection.UUID)
	m.Name = types.StringValue(policyCollection.Name)
	m.Description = types.StringPointerValue(policyCollection.Description)
	m.CloudProvider = types.StringValue(string(policyCollection.Provider))
	m.AutoUpdate = types.BoolValue(policyCollection.AutoUpdate)
	m.System = types.BoolValue(policyCollection.System)
	m.Dynamic = types.BoolValue(policyCollection.IsDynamic)

	dynamicConfig, d := typehelpers.ObjectValue(
		ctx,
		policyCollection.RepositoryView,
		func() (*PolicyCollectionDynamicConfig, diag.Diagnostics) {
			return &PolicyCollectionDynamicConfig{
				RepositoryUUID:     types.StringValue(*policyCollection.RepositoryConfig.UUID),
				Namespace:          types.StringValue(policyCollection.RepositoryView.Namespace),
				BranchName:         types.StringValue(policyCollection.RepositoryView.BranchName),
				PolicyDirectories:  typehelpers.StringsList(policyCollection.RepositoryView.PolicyDirectories),
				PolicyFileSuffixes: typehelpers.StringsList(policyCollection.RepositoryView.PolicyFileSuffix),
			}, nil
		},
	)
	errors.AddAttributeDiags(&diags, d, "dynamic_config")
	m.DynamicConfig = dynamicConfig
	m.RoleAssignmentTarget = types.StringValue("policy-collection:" + policyCollection.UUID)

	return diags
}

// PolicyCollectionDatasource is the model for a policy collection data source.
type PolicyCollectionDataSource struct {
	PolicyCollectionResource
}

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
