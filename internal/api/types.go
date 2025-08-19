// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringEnum represents a generic string-based enum type.
type StringEnum string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (se *StringEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*se = StringEnum(s)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (se StringEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(se))
}

// String implements the fmt.Stringer interface.
func (se StringEnum) String() string {
	return string(se)
}

// StringsList converts a types.List to a list of strings.
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
