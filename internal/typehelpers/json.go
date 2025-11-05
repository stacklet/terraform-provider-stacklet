// Copyright (c) 2025 - Stacklet, Inc.

package typehelpers

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// JSONString is a string containing JSON which gets normalized for sorting/whitespace.
func JSONString(s *string) (types.String, error) {
	if s == nil {
		return types.StringNull(), nil
	}

	var data any
	if err := json.Unmarshal([]byte(*s), &data); err != nil {
		return types.StringNull(), err
	}
	newString, err := json.Marshal(data)
	if err != nil {
		return types.StringNull(), err
	}
	return types.StringValue(string(newString)), nil
}
