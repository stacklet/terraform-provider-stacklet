package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyCollectionItemResource is the model for a policy collection item resource.
type PolicyCollectionItemResource struct {
	ID             types.String `tfsdk:"id"`
	CollectionUUID types.String `tfsdk:"collection_uuid"`
	PolicyUUID     types.String `tfsdk:"policy_uuid"`
	PolicyVersion  types.Int32  `tfsdk:"policy_version"`
}
