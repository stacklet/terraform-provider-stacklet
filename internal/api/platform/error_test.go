// Copyright Stacklet, Inc. 2025, 2026

package platform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

func TestFromProblems(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		problems []problem
		want     error
	}{
		{"none", nil, nil},
		{"not found", []problem{{Kind: "NotFound", Message: "gone"}}, api.NotFound{Message: "gone"}},
		{"conflict", []problem{{Kind: "Conflict", Message: "try again"}}, api.Conflict{Message: "try again"}},
		{"other", []problem{{Kind: "InvalidInput", Message: "bad key"}}, api.APIError{Kind: "InvalidInput", Detail: "bad key"}},
		{
			"first of several",
			[]problem{{Kind: "InvalidInput", Message: "first"}, {Kind: "NotFound", Message: "second"}},
			api.APIError{Kind: "InvalidInput", Detail: "first"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fromProblems(ctx, tt.problems))
		})
	}
}
