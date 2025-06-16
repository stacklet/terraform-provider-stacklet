// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ObjectMap returns a basetypes.MapValue from a list of objects.
func ObjectMapFromList[ElemType WithAttributes, ItemType any](l []ItemType, buildElement func(ItemType) (string, map[string]attr.Value, diag.Diagnostics)) (basetypes.MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var emptyElem ElemType
	attrTypes := emptyElem.AttributeTypes()
	elemType := types.ObjectType{AttrTypes: attrTypes}

	if l == nil {
		return types.MapNull(elemType), diags
	}

	m := make(map[string]attr.Value)
	for _, entry := range l {
		key, values, d := buildElement(entry)
		diags.Append(d...)
		if diags.HasError() {
			break
		}

		v, d := types.ObjectValue(attrTypes, values)
		diags.Append(d...)
		if diags.HasError() {
			break
		}
		m[key] = v
	}

	result, d := types.MapValue(elemType, m)
	diags.Append(d...)
	return result, diags
}
