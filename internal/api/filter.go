// Copyright (c) 2025 - Stacklet, Inc.

package api

// FilterElementinput define an element filter input.
type FilterElementInput struct {
	Single FilterValueInput `json:"single"`
}

// FilterValueInput is a filter for a single value.
type FilterValueInput struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// NewFilterElementSingle returns a populated FilterElementSingle.
func NewFieldMatchFilter(name string, value any) FilterElementInput {
	return FilterElementInput{
		Single: FilterValueInput{
			Name:     name,
			Operator: "equals",
			Value:    value,
		},
	}
}
