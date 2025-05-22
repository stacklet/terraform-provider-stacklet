package types

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NullableString returns the proper type for a nullable string.
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
