package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyDataSource is the model for policy data sources
type PolicyDataSource struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Version       types.Number `tfsdk:"version"`
}
