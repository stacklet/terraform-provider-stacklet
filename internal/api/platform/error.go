// Copyright Stacklet, Inc. 2025, 2026

package platform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// fromProblems returns an error from a list of API problems.
func fromProblems(ctx context.Context, problems []problem) error {
	if len(problems) == 0 {
		return nil
	}
	for _, problem := range problems[1:] {
		info := map[string]any{"kind": problem.Kind, "message": problem.Message}
		tflog.Error(ctx, "discarding additional error", info)
	}
	switch problems[0].Kind {
	case "NotFound":
		return api.NotFound{Message: problems[0].Message}
	case "Conflict":
		return api.Conflict{Message: problems[0].Message}
	}
	return api.APIError{Kind: problems[0].Kind, Detail: problems[0].Message}
}

// problem contains the details for an API query error.
type problem struct {
	Kind    string `graphql:"__typename"`
	Message string
}
