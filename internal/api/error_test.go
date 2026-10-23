// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIError_URLError(t *testing.T) {
	// http.Client.Do wraps transport errors in *url.Error. Without unwrapping,
	// the diagnostic would show `Post "https://...": field not found` instead of
	// the clean inner message.
	inner := errors.New("field not found; invalid type")
	urlErr := &url.Error{Op: "Post", URL: "https://api.example.com/graphql", Err: inner}
	err := newAPIError(urlErr)

	assert.Equal(t, "API Error", err.Summary())
	assert.Equal(t, "field not found; invalid type", err.Error())
}

func TestNewAPIError_PlainError(t *testing.T) {
	err := newAPIError(errors.New("something went wrong"))

	assert.Equal(t, "API Error", err.Summary())
	assert.Equal(t, "something went wrong", err.Error())
}

func TestFromProblems(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		problems []problem
		want     error
	}{
		{"none", nil, nil},
		{"not found", []problem{{Kind: "NotFound", Message: "gone"}}, NotFound{"gone"}},
		{"conflict", []problem{{Kind: "Conflict", Message: "try again"}}, Conflict{"try again"}},
		{"other", []problem{{Kind: "InvalidInput", Message: "bad key"}}, apiError{Kind: "InvalidInput", Detail: "bad key"}},
		{
			"first of several",
			[]problem{{Kind: "InvalidInput", Message: "first"}, {Kind: "NotFound", Message: "second"}},
			apiError{Kind: "InvalidInput", Detail: "first"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fromProblems(ctx, tt.problems))
		})
	}
}
