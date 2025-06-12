// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
