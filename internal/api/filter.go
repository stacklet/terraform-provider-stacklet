// Copyright (c) 2025 - Stacklet, Inc.

package api

// FilterElementInput define an element filter input.
// Matches the platform API structure from:
// https://github.com/stacklet/platform/blob/main/src/stacklet/platform/filters/input.py
type FilterElementInput struct {
	Single *FilterValueInput `json:"single,omitempty"`
}

// FilterValueInput is a filter for a single value.
type FilterValueInput struct {
	Name     string `json:"name"`
	Operator string `json:"operator,omitempty"`
	Value    any    `json:"value"`
}

// newExactMatchFilter returns a populated FilterElementInput for an exact match.
func newExactMatchFilter(name string, value any) FilterElementInput {
	return FilterElementInput{
		Single: &FilterValueInput{
			Name:     name,
			Operator: "equals",
			Value:    value,
		},
	}
}

// newSimpleFilter returns a populated FilterElementInput for a simple filter without an operator.
// Used for filters that don't require an operator (like user username filtering).
func newSimpleFilter(name string, value any) FilterElementInput {
	return FilterElementInput{
		Single: &FilterValueInput{
			Name:  name,
			Value: value,
		},
	}
}
