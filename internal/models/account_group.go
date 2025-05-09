package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccountGroupResource is the model for account group resources
type AccountGroupResource struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Regions       types.List   `tfsdk:"regions"`
}

// AccountGroupDataSource is the model for account group data sources
type AccountGroupDataSource AccountGroupResource
