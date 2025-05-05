package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// nullableString return the proper type for a nullable string
func nullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
