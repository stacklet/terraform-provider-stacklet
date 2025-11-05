// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type WithAttributes interface {
	AttributeTypes() map[string]attr.Type
}

// ObjectValue returns a basetypes.ObjectValue from a type.
func ObjectValue[Type WithAttributes, Value any](ctx context.Context, v *Value, construct func() (*Type, diag.Diagnostics)) (basetypes.ObjectValue, diag.Diagnostics) {
	var empty Type
	var diags diag.Diagnostics
	if v == nil {
		return NullObject(empty), diags
	}
	objPtr, d := construct()
	diags.Append(d...)
	if diags.HasError() || objPtr == nil {
		return NullObject(empty), diags
	}
	return types.ObjectValueFrom(ctx, empty.AttributeTypes(), *objPtr)
}

// NullObject returns a basetype.ObjectValue for the provided type with empty values.
func NullObject(t WithAttributes) basetypes.ObjectValue {
	return basetypes.NewObjectNull(t.AttributeTypes())
}

// FilteredObject returns an object with only the specified attributes.
func FilteredObject[Type WithAttributes](obj basetypes.ObjectValue, fields []string) (basetypes.ObjectValue, diag.Diagnostics) {
	var empty Type

	if obj.IsNull() || obj.IsUnknown() {
		return NullObject(empty), nil
	}

	attrs := obj.Attributes()
	values := make(map[string]attr.Value)
	for _, field := range fields {
		values[field] = attrs[field]
	}
	return types.ObjectValue(empty.AttributeTypes(), values)
}

// UpdateObject returns an object with the specified attributes added/updated.
func UpdatedObject(ctx context.Context, obj basetypes.ObjectValue, attrs map[string]attr.Value) (basetypes.ObjectValue, diag.Diagnostics) {
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
	id, _ := attr.(basetypes.StringValue)
	return id.ValueString()
}
