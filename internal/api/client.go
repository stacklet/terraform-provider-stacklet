// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hasura/go-graphql-client"
)

// NewClient returns a configured graphql Client.
func NewClient(ctx context.Context, endpoint string, apiKey string) *graphql.Client {
	tfLog := hclog.LevelFromString(os.Getenv("TF_LOG"))
	logBody := tfLog == hclog.Debug || tfLog == hclog.Trace

	httpClient := &http.Client{
		Transport: &authTransport{
			APIKey: apiKey,
			Base: &logTransport{
				Ctx:     ctx,
				Base:    http.DefaultTransport,
				LogBody: logBody,
			},
		},
	}
	return graphql.NewClient(endpoint, httpClient)
}

// authTransport is an http.Transport that adds authorization header.
type authTransport struct {
	APIKey string
	Base   http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.APIKey)
	return t.Base.RoundTrip(req)
}

// logTransport is an http.Transport that logs requests/responses.
type logTransport struct {
	Ctx     context.Context
	Base    http.RoundTripper
	LogBody bool
}

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody, respBody string
	var err error

	// Decode the request before performing it since it will otherwise consume the body
	if t.LogBody {
		reqBody, err = decodeRequestBody(req)
		if err != nil {
			return nil, err
		}
	}

	tflog.Debug(
		t.Ctx,
		"Performing GraphQL request",
		map[string]any{
			"req_method": req.Method,
			"req_url":    req.URL.String(),
			"req_body":   reqBody,
		},
	)

	resp, reqErr := t.Base.RoundTrip(req)

	if t.LogBody {
		respBody, err = decodeResponseBody(resp)
		if err != nil {
			return nil, err
		}
	}

	tflog.Debug(
		t.Ctx,
		"Got GraphQL response",
		map[string]any{
			"req_method":       req.Method,
			"req_url":          req.URL.String(),
			"resp_status":      resp.Status,
			"resp_status_code": resp.StatusCode,
			"req_body":         reqBody,
			"resp_body":        respBody,
		},
	)

	return resp, reqErr
}

func decodeRequestBody(req *http.Request) (string, error) {
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

func decodeResponseBody(resp *http.Response) (string, error) {
	if resp.Body == nil || resp.Body == http.NoBody {
		return "", nil
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body = io.NopCloser(bytes.NewReader(content))

	return string(content), nil
}
