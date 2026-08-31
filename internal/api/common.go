// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TerraformModule is the data returned for terraform module definitions.
type TerraformModule struct {
	RepositoryURL string  `graphql:"repositoryURL"`
	Source        string  `graphql:"source"`
	Version       *string `graphql:"version"`
	VariablesJSON string  `graphql:"variablesJSON"`
}

// StatusInfo is the data returned for the status of an entity's background
// operations, as a summary over the individual outcomes.
type StatusInfo struct {
	Status  string         `graphql:"status"`
	Details []StatusDetail `graphql:"details"`
}

// StatusDetail is the outcome of one named operation's most recent completed
// attempt, plus whether a new attempt is underway.
//
// Status is nullable: an operation that has never completed one has no outcome
// to report, which is not a failure. Timestamps are kept as the strings the API
// returns, since they are only ever surfaced for a human to read.
type StatusDetail struct {
	Component    string  `graphql:"component"`
	Status       *string `graphql:"status"`
	At           *string `graphql:"at"`
	Message      *string `graphql:"message"`
	RunningSince *string `graphql:"runningSince"`
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

	tagsMap := make(map[string]attr.Value)
	for _, tag := range t {
		tagsMap[tag.Key] = types.StringValue(tag.Value)
	}

	result, d := types.MapValue(types.StringType, tagsMap)
	diags.Append(d...)
	return result, diags
}
