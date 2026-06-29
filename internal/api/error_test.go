// Copyright Stacklet, Inc. 2025, 2026

package api

import (
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
