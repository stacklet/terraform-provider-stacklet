// Copyright Stacklet, Inc. 2025, 2026

// package api provides functionality shared by the API clients.
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// NewTransport returns an http.RoundTripper which authenticates and logs requests.
func NewTransport(apiKey string, version string) http.RoundTripper {
	tfLog := hclog.LevelFromString(os.Getenv("TF_LOG"))
	logBody := tfLog == hclog.Debug || tfLog == hclog.Trace

	return &AuthTransport{
		APIKey:  apiKey,
		Version: version,
		Base: &LogTransport{
			Base:    http.DefaultTransport,
			LogBody: logBody,
		},
	}
}

// AuthTransport is an http.Transport that adds authorization header.
type AuthTransport struct {
	APIKey  string
	Version string
	Base    http.RoundTripper
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.APIKey == "" {
		return nil, fmt.Errorf("no API key provided")
	}
	req.Header.Set("Authorization", "Bearer "+t.APIKey)
	req.Header.Set("User-Agent", "terraform-provider-stacklet/"+t.Version)
	return t.Base.RoundTrip(req)
}

// LogTransport is an http.Transport that logs requests/responses through the
// request's context, so the logs apply whatever masking the caller set up.
type LogTransport struct {
	Base    http.RoundTripper
	LogBody bool
}

func (t *LogTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody, respBody string
	var err error

	// Decode the request before performing it since it will otherwise consume the body
	if t.LogBody {
		reqBody, err = DecodeRequestBody(req)
		if err != nil {
			return nil, err
		}
	}

	tflog.Debug(
		req.Context(),
		"Performing API request",
		map[string]any{
			"req_method": req.Method,
			"req_url":    req.URL.String(),
			"req_body":   reqBody,
		},
	)

	resp, reqErr := t.Base.RoundTrip(req)

	if resp != nil {
		if t.LogBody {
			respBody, err = DecodeResponseBody(resp)
			if err != nil {
				return nil, err
			}
		}

		tflog.Debug(
			req.Context(),
			"Got API response",
			map[string]any{
				"req_method":       req.Method,
				"req_url":          req.URL.String(),
				"resp_status":      resp.Status,
				"resp_status_code": resp.StatusCode,
				"req_body":         reqBody,
				"resp_body":        respBody,
			},
		)
	}

	return resp, reqErr
}

// DecodeRequestBody returns the request body, leaving it readable again.
func DecodeRequestBody(req *http.Request) (string, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return "", nil
	}

	content, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewReader(content))

	return string(content), nil
}

// DecodeResponseBody returns the response body, leaving it readable again.
func DecodeResponseBody(resp *http.Response) (string, error) {
	if resp.Body == nil || resp.Body == http.NoBody {
		return "", nil
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body = io.NopCloser(bytes.NewReader(content))

	return string(content), nil
}
