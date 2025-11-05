// Copyright (c) 2025 - Stacklet, Inc.

package typehelpers

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringsList returns a list of values of string type.
func StringsList(l []string) types.List {
	sl := make([]attr.Value, len(l))
	for i, item := range l {
		sl[i] = types.StringValue(item)
	}
	lv, _ := types.ListValue(types.StringType, sl)
	return lv
}

// ObjectList returns a types.List from a list of objects.
func ObjectList[ElemType WithAttributes, ItemType any](l []ItemType, buildElement func(ItemType) (map[string]attr.Value, diag.Diagnostics)) (types.List, diag.Diagnostics) {
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

// ListItemsIdentifiers returns a list of identifiers from the specified attribute from a types.List of objects.
func ListItemsIdentifiers(l types.List, attrName string) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}

	elems := l.Elements()
	ids := make([]string, len(elems))
	for i, elem := range elems {
		if obj, ok := elem.(types.Object); ok {
			ids[i] = ObjectStringIdentifier(obj, attrName)
		}
	}
	return ids
}

// ListSortedEntries returns a list with same entries as the original one, but sorted by the specified string attribute according to the provided values, with extra entries at the end.
func ListSortedEntries[Type WithAttributes](l types.List, attrName string, attrValues []string) (types.List, diag.Diagnostics) {
	var empty Type
	attrTypes := empty.AttributeTypes()

	if l.IsNull() || l.IsUnknown() {
		types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	if attrValues == nil {
		return l, nil
	}

	listIdentifiers := ListItemsIdentifiers(l, attrName)

	elems := l.Elements()

	elemsByID := make(map[string]attr.Value)
	for i, id := range listIdentifiers {
		elemsByID[id] = elems[i]
	}

	values := []attr.Value{}
	seen := make(map[string]bool)
	for _, id := range attrValues {
		if elem, ok := elemsByID[id]; ok {
			values = append(values, elem)
			seen[id] = true
		}
	}

	// add extra elements at the end
	for i, id := range listIdentifiers {
		if !seen[id] {
			values = append(values, elems[i])
		}
	}
	return types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
}
