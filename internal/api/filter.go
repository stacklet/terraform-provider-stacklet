// Copyright Stacklet, Inc. 2025, 2026

package api

// filterElementInput define an element filter input.
// Matches the platform API structure from:
// https://github.com/stacklet/platform/blob/main/src/stacklet/platform/filters/input.py
type filterElementInput struct {
	Single *filterValueInput `json:"single,omitempty"`
}

func (i filterElementInput) GetGraphQLType() string {
	return "FilterElementInput"
}

// filterValueInput is a filter for a single value.
type filterValueInput struct {
	Name     string `json:"name"`
	Operator string `json:"operator,omitempty"`
	Value    any    `json:"value"`
}

// newExactMatchFilter returns a populated filterElementInput for an exact match.
func newExactMatchFilter(name string, value any) filterElementInput {
	return filterElementInput{
		Single: &filterValueInput{
			Name:     name,
			Operator: "equals",
			Value:    value,
		},
	}
}

// newSimpleFilter returns a populated filterElementInput for a simple filter without an operator.
// Used for filters that don't require an operator (like user username filtering).
func newSimpleFilter(name string, value any) filterElementInput {
	return filterElementInput{
		Single: &filterValueInput{
			Name:  name,
			Value: value,
		},
	}
}
