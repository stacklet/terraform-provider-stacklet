package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullableString return the proper type for a nullable string
func NullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
