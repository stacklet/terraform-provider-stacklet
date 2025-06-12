// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// StringsList returns a list of values of string type.
func StringsList(l []string) basetypes.ListValue {
	sl := make([]attr.Value, len(l))
	for i, item := range l {
		sl[i] = types.StringValue(item)
	}
	lv, _ := types.ListValue(types.StringType, sl)
	return lv
}
