// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// ObjectList returns a basetypes.ListValue from a list of objects.
func ObjectList[ElemType WithAttributes, ItemType any](l []ItemType, buildElement func(ItemType) (map[string]attr.Value, diag.Diagnostics)) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var emptyElem ElemType
	attrTypes := emptyElem.AttributeTypes()
	elemType := types.ObjectType{AttrTypes: attrTypes}

	if l == nil {
		return types.ListNull(elemType), diags
	}

	values := []attr.Value{}
	for _, entry := range l {
		objValues, d := buildElement(entry)
		diags.Append(d...)
		if diags.HasError() {
			break
		}

		v, d := types.ObjectValue(attrTypes, objValues)
		diags.Append(d...)
		if diags.HasError() {
			break
		}
		values = append(values, v)
	}

	result, d := types.ListValue(elemType, values)
	diags.Append(d...)
	return result, diags
}
