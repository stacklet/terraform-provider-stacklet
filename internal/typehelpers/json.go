// Copyright (c) 2025 - Stacklet, Inc.

package typehelpers

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

// JSONString returns a string normalized for sorting/whitespace.
func JSONString(s *string) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	if s == nil {
		return types.StringNull(), diags
	}

	var data any
	if err := json.Unmarshal([]byte(*s), &data); err != nil {
		errors.AddDiagError(&diags, err)
		return types.StringNull(), diags
	}
	newString, err := json.Marshal(data)
	if err != nil {
		errors.AddDiagError(&diags, err)
		return types.StringNull(), diags
	}
	return types.StringValue(string(newString)), diags
}
