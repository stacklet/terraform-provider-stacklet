// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	return basetypes.NewObjectValueFrom(ctx, empty.AttributeTypes(), *objPtr)
}

// NullObject returns a basetype.ObjectValue for the provided type with empty values.
func NullObject(t WithAttributes) basetypes.ObjectValue {
	return basetypes.NewObjectNull(t.AttributeTypes())
}
