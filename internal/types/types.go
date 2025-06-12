// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// NullableBool returns the proper type for a nullable boolean.
func NullableBool(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

// NullableInt returns the proper type for a nullable integer.
func NullableInt(i *int) types.Int32 {
	if i == nil {
		return types.Int32Null()
	}
	return types.Int32Value(int32(*i))
}

// NullableFloat returns the proper type for a nullable float.
func NullableFloat(f *float32) types.Float32 {
	if f == nil {
		return types.Float32Null()
	}
	return types.Float32Value(*f)
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
