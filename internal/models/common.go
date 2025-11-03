// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TerraformModule is the model for terraform modules definitions.
type TerraformModule struct {
	RepositoryURL types.String `tfsdk:"repository_url"`
	Source        types.String `tfsdk:"source"`
	VariablesJSON types.String `tfsdk:"variables_json"`
	Version       types.String `tfsdk:"version"`
}

func (c TerraformModule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository_url": types.StringType,
		"source":         types.StringType,
		"variables_json": types.StringType,
		"version":        types.StringType,
	}
}

// Tag is the model for tags.
type Tag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (t Tag) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}
