// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ListItemsIdentifiers returns a list of identifiers from the specified attribute from a types.List of objects.
func ListItemsIdentifiers(l types.List, attrName string) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}

	elems := l.Elements()
	ids := make([]string, len(elems))
	for i, elem := range elems {
		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		attr := obj.Attributes()[attrName]
		id, ok := attr.(basetypes.StringValue)
		if !ok {
			continue
		}
		ids[i] = id.ValueString()
	}
	return ids
}

// ItemsByIdentifier returns a map of objects by the specified attribute from a types.List of Objects.
func ItemsByIdentifier(l types.List, attrName string) map[string]types.Object {
	items := map[string]types.Object{}

	for _, elem := range l.Elements() {
		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		attr := obj.Attributes()[attrName]
		id, ok := attr.(basetypes.StringValue)
		if !ok {
			continue
		}
		items[id.ValueString()] = obj
	}

	return items
}
