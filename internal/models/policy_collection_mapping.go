// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// PolicyCollectionMappingResource is the model for a policy collection mapping resource.
type PolicyCollectionMappingResource struct {
	ID             types.String `tfsdk:"id"`
	CollectionUUID types.String `tfsdk:"collection_uuid"`
	PolicyUUID     types.String `tfsdk:"policy_uuid"`
	PolicyVersion  types.Int32  `tfsdk:"policy_version"`
}

func (m *PolicyCollectionMappingResource) Update(policyCollectionMapping *api.PolicyCollectionMapping) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(policyCollectionMapping.ID)
	m.CollectionUUID = types.StringValue(policyCollectionMapping.Collection.UUID)
	m.PolicyUUID = types.StringValue(policyCollectionMapping.Policy.UUID)
	m.PolicyVersion = types.Int32Value(int32(policyCollectionMapping.Policy.Version))

	return diags
}
