// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"errors"
	"net/url"
)

// APIError represent an error interacting with the API.
type APIError struct {
	Kind   string
	Detail string
}

// Summary returns the error summary message.
func (e APIError) Summary() string {
	return e.Kind
}

// Error returns the error message.
func (e APIError) Error() string {
	return e.Detail
}

// NewAPIError returns an APIError from an error.
func NewAPIError(err error) APIError {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return NewAPIError(urlErr.Err)
	}
	return APIError{Kind: "API Error", Detail: err.Error()}
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

// Conflict represents an error raised when a change isn't possible, or wasn't
// completed, because of the state of some other part of the system.
type Conflict struct {
	Message string
}

// Summary returns the error summary.
func (e Conflict) Summary() string {
	return "Conflict"
}

// Error returns the error message.
func (e Conflict) Error() string {
	return e.Message
}
