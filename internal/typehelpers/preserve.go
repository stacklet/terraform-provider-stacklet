// Copyright Stacklet, Inc. 2025, 2026

package typehelpers

import "github.com/hashicorp/terraform-plugin-framework/types"

// PreserveIfNull returns original when current is null but original is not.
// Use when the API omits a field for its zero value, which would cause a perpetual plan diff.
func PreserveIfNull[T interface{ IsNull() bool }](current, original T) T {
	if current.IsNull() && !original.IsNull() {
		return original
	}
	return current
}

// PreserveIfEmptyJSON returns original when current serialises as "{}".
// Use when the API returns an empty dictionary for both null and empty values.
func PreserveIfEmptyJSON(current, original types.String) types.String {
	if current.ValueString() == "{}" {
		return original
	}
	return current
}
