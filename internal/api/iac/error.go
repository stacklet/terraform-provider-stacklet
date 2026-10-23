// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// errorPayload is the body of an error response.
//
// The API returns a "detail" field which is either a string or a list of
// validation errors. Errors returned by the API gateway have a "message" field instead.
type errorPayload struct {
	Detail  json.RawMessage `json:"detail"`
	Message string          `json:"message"`
}

type validationError struct {
	Loc []any  `json:"loc"`
	Msg string `json:"msg"`
}

// newResponseError returns an error for a non-successful API response.
func newResponseError(statusCode int, content []byte) error {
	detail := errorDetail(content)
	if detail == "" {
		detail = http.StatusText(statusCode)
	}

	switch statusCode {
	case http.StatusNotFound:
		return api.NotFound{Message: detail}
	case http.StatusUnauthorized, http.StatusForbidden:
		return api.APIError{Kind: "Authorization Error", Detail: detail}
	default:
		return api.APIError{Kind: "API Error", Detail: detail}
	}
}

func errorDetail(content []byte) string {
	var payload errorPayload
	if err := json.Unmarshal(content, &payload); err != nil {
		return strings.TrimSpace(string(content))
	}

	var detail string
	if err := json.Unmarshal(payload.Detail, &detail); err == nil {
		return detail
	}

	var validationErrors []validationError
	if err := json.Unmarshal(payload.Detail, &validationErrors); err == nil && len(validationErrors) > 0 {
		msgs := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			msgs[i] = e.Msg
			if len(e.Loc) > 0 {
				loc := make([]string, len(e.Loc))
				for j, l := range e.Loc {
					loc[j] = fmt.Sprint(l)
				}
				msgs[i] = strings.Join(loc, ".") + ": " + e.Msg
			}
		}
		return strings.Join(msgs, "\n")
	}

	return payload.Message
}
