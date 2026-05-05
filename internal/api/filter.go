// Copyright Stacklet, Inc. 2025, 2026

package api

// filterBooleanOperator represents a boolean operation fora filter.
type filterBooleanOperator StringEnum

const (
	filterBooleanAND = filterBooleanOperator("AND")
	filterBooleanOR  = filterBooleanOperator("OR")
	filterBooleanNOT = filterBooleanOperator("NOT")
)

// filterElementInput define an element filter input.
// Matches the platform API structure from:
// https://github.com/stacklet/platform/blob/main/src/stacklet/platform/filters/input.py
type filterElementInput struct {
	Single   *filterValueInput            `json:"single,omitempty"`
	Multiple *filterBooleanOperationInput `json:"multiple,omitempty"`
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

// filterBooleanOperationInput is a combined filter.
type filterBooleanOperationInput struct {
	Operands []filterElementInput  `json:"operands"`
	Operator filterBooleanOperator `json:"operator"`
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

// newCompositeFilter returns a populated filterElementInput for a composite
// filter with multiple conditions with a single operator.
func newCompositeFilter(filters []filterElementInput, operator filterBooleanOperator) filterElementInput {
	return filterElementInput{
		Multiple: &filterBooleanOperationInput{
			Operands: filters,
			Operator: operator,
		},
	}
}
