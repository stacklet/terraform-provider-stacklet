// Copyright (c) 2025 - Stacklet, Inc.

package typehelpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WithAttributes interface {
	AttributeTypes() map[string]attr.Type
}

// ObjectValue returns a types.Object from a type.
func ObjectValue[Type WithAttributes, Value any](ctx context.Context, v *Value, construct func() (*Type, diag.Diagnostics)) (types.Object, diag.Diagnostics) {
	var empty Type
	attrTypes := empty.AttributeTypes()
	nullObj := types.ObjectNull(attrTypes)

	if v == nil {
		return nullObj, nil
	}
	objPtr, diags := construct()
	if diags.HasError() || objPtr == nil {
		return nullObj, diags
	}
	return types.ObjectValueFrom(ctx, attrTypes, *objPtr)
}

// FilteredObject returns an object with only the specified attributes.
func FilteredObject[Type WithAttributes](obj types.Object, fields []string) (types.Object, diag.Diagnostics) {
	var empty Type
	attrTypes := empty.AttributeTypes()

	if obj.IsNull() || obj.IsUnknown() {
		return types.ObjectNull(attrTypes), nil
	}

	attrs := obj.Attributes()
	values := make(map[string]attr.Value)
	for _, field := range fields {
		values[field] = attrs[field]
	}
	return types.ObjectValue(attrTypes, values)
}

// UpdateObject returns an object with the specified attributes added/updated.
func UpdatedObject(ctx context.Context, obj types.Object, attrs map[string]attr.Value) (types.Object, diag.Diagnostics) {
	objAttrs := obj.Attributes()
	objTypes := obj.AttributeTypes(ctx)
	for name, value := range attrs {
		objAttrs[name] = value
		objTypes[name] = value.Type(ctx)
	}
	return types.ObjectValue(objTypes, objAttrs)
}

// ObjectStringIdentifier returns the string identifier for the object, from the provided attribute (which must be of type.String).
func ObjectStringIdentifier(obj types.Object, name string) string {
	attr := obj.Attributes()[name]
	id, _ := attr.(types.String)
	return id.ValueString()
}
