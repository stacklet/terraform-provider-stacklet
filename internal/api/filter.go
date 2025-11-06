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

// newExactMatchFilter returns a populated FilterElementInput for an exact match.
func newExactMatchFilter(name string, value any) FilterElementInput {
	return FilterElementInput{
		Single: FilterValueInput{
			Name:     name,
			Operator: "equals",
			Value:    value,
		},
	}
}
