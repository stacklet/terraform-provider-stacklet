// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// apiError represent an error interacting with the API.
type apiError struct {
	Kind   string
	Detail string
}

// Error returns the error summary message.
func (e apiError) Summary() string {
	return e.Kind
}

// Error returns the error message.
func (e apiError) Error() string {
	return e.Detail
}

// newAPIError returns an apiError from an error.
func newAPIError(err error) apiError {
	return apiError{Kind: "API Error", Detail: err.Error()}
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

// fromProblems returns an error from a list of API problems.
func fromProblems(ctx context.Context, problems []problem) error {
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
	return apiError{Kind: problems[0].Kind, Detail: problems[0].Message}
}

// problem contains the details for an API query error.
type problem struct {
	Kind    string `graphql:"__typename"`
	Message string
}
