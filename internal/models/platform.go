// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PlatformDataSource is the model for the platform data source.
type PlatformDataSource struct {
	ID               types.String `tfsdk:"id"`
	ExternalID       types.String `tfsdk:"external_id"`
	ExecutionRegions types.List   `tfsdk:"execution_regions"`
	DefaultRole      types.String `tfsdk:"default_role"`
}
