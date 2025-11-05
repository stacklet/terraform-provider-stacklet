// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// NewTagsList returns a TagsList from a map. Elements are sorted by key.
func NewTagsList(tags types.Map) TagsList {
	tagsList := make(TagsList, 0)

	if tags.IsNull() || tags.IsUnknown() {
		return tagsList
	}

	tagsMap := tags.Elements()
	keys := make([]string, 0, len(tagsMap))
	for key := range tagsMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if strVal, ok := tagsMap[key].(types.String); ok {
			tagsList = append(tagsList, Tag{
				Key:   key,
				Value: strVal.ValueString(),
			})
		}
	}

	return tagsList
}

// TagsMap converts a list of tags to a map of key-value pairs.
func (t TagsList) TagsMap() (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	tagsMap := map[string]attr.Value{}
	for _, tag := range t {
		tagsMap[tag.Key] = types.StringValue(tag.Value)
	}

	result, d := types.MapValue(types.StringType, tagsMap)
	diags.Append(d...)
	return result, diags
}
