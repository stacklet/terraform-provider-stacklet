// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyDataSource is the model for policy data sources.
type PolicyDataSource struct {
	ID              types.String `tfsdk:"id"`
	UUID            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	CloudProvider   types.String `tfsdk:"cloud_provider"`
	Version         types.Int32  `tfsdk:"version"`
	Category        types.List   `tfsdk:"category"`
	Mode            types.String `tfsdk:"mode"`
	ResourceType    types.String `tfsdk:"resource_type"`
	Path            types.String `tfsdk:"path"`
	SourceJSON      types.String `tfsdk:"source_json"`
	SourceYAML      types.String `tfsdk:"source_yaml"`
	System          types.Bool   `tfsdk:"system"`
	UnqualifiedName types.String `tfsdk:"unqualified_name"`
}
