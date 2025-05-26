package api

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// APIError represent an error interacting with the API.
type APIError struct {
	Summary string
	Detail  string
}

// Error returns the error message.
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Summary, e.Detail)
}

// NotFound represents an error raised when an API resource is not found.
type NotFound struct {
	Message string
}

// Summary returns the error summary.
func (e NotFound) Summary() string {
	return "Not Found"
}

// Error returns the error message.
func (e NotFound) Error() string {
	return e.Message
}

func FromProblems(ctx context.Context, problems []Problem) error {
	if len(problems) == 0 {
		return nil
	}
	for _, problem := range problems[1:] {
		info := map[string]any{"kind": problem.Kind, "message": problem.Message}
		tflog.Error(ctx, "discarding additional error", info)
	}
	if problems[0].Kind == "NotFound" {
		return NotFound{problems[0].Message}
	}
	return APIError{problems[0].Kind, problems[0].Message}
}

type Problem struct {
	Kind    string `graphql:"__typename"`
	Message string
}
