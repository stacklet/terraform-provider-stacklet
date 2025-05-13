package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccountGroupItemResource is the model for an account group item resource.
type AccountGroupItemResource struct {
	ID            types.String `tfsdk:"id"`
	GroupUUID     types.String `tfsdk:"group_uuid"`
	AccountKey    types.String `tfsdk:"account_key"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
}
