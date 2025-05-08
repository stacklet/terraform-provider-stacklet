package api

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullableString converts a types.String to a string pointer that can be null
func NullableString(s types.String) *string {
	if s.IsNull() {
		return nil
	}

	str := s.ValueString()
	return &str
}
