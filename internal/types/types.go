package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NullableString return the proper type for a nullable string.
func NullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

// StringsList returns a list of values of string type.
func StringsList(l []string) basetypes.ListValue {
	sl := make([]attr.Value, len(l))
	for i, item := range l {
		sl[i] = types.StringValue(item)
	}
	lv, _ := types.ListValue(types.StringType, sl)
	return lv
}
