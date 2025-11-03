// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// TerraformModule is the data returned for terraform module definitions.
type TerraformModule struct {
	RepositoryURL string `graphql:"repositoryURL"`
	Source        string
	Version       *string
	VariablesJSON string `graphql:"variablesJSON"`
}

// Tag is the data for a tag.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TagsList is a list of tags.
type TagsList []Tag

// TagsMap converts a list of tags to a map of key-value pairs.
func (t TagsList) TagsMap() (basetypes.MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tagsMap := map[string]attr.Value{}
	for _, tag := range t {
		tagsMap[tag.Key] = types.StringValue(tag.Value)
	}

	result, d := types.MapValue(types.StringType, tagsMap)
	diags.Append(d...)
	return result, diags
}
