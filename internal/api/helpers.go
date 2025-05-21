package api

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullableString converts a types.String to a string pointer which can be null.
func NullableString(s types.String) *string {
	if s.IsNull() {
		return nil
	}

	str := s.ValueString()
	return &str
}

// NullableBool converts a types.Bool to a bool pointer which can be null.
func NullableBool(b types.Bool) *bool {
	if b.IsNull() {
		return nil
	}

	bv := b.ValueBool()
	return &bv
}

// StringsList concerts a types.List to a list of strings.
func StringsList(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	elements := l.Elements()
	sl := make([]string, len(elements))
	for i, element := range elements {
		str, _ := element.(types.String)
		sl[i] = str.ValueString()
	}
	return sl
}
