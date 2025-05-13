package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccountGroupMappingResource is the model for an account group mapping resource.
type AccountGroupMappingResource struct {
	ID            types.String `tfsdk:"id"`
	GroupUUID     types.String `tfsdk:"group_uuid"`
	AccountKey    types.String `tfsdk:"account_key"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
}
