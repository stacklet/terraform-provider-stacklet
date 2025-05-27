// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyCollectionMappingResource is the model for a policy collection mapping resource.
type PolicyCollectionMappingResource struct {
	ID             types.String `tfsdk:"id"`
	CollectionUUID types.String `tfsdk:"collection_uuid"`
	PolicyUUID     types.String `tfsdk:"policy_uuid"`
	PolicyVersion  types.Int32  `tfsdk:"policy_version"`
}
