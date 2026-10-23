// Copyright Stacklet, Inc. 2025, 2026

// package iac provides access to the IaC (Sinistral) REST API.
package iac

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// ClientConfig is the configuration for the API client.
type ClientConfig struct {
	Endpoint string
	APIKey   string
	Version  string
}

// client is a JSON client for the REST API.
type client struct {
	httpClient *http.Client
	endpoint   string
}

// newClient returns a configured REST client.
func newClient(config ClientConfig) *client {
	return &client{
		httpClient: &http.Client{
			Transport: api.NewTransport(config.APIKey, config.Version),
		},
		endpoint: config.Endpoint,
	}
}

// Get performs a GET request, decoding the result into response.
func (c *client) Get(ctx context.Context, path string, query url.Values, response any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, response)
}

// Post performs a POST request with payload as body, decoding the result into response.
func (c *client) Post(ctx context.Context, path string, payload any, response any) error {
	return c.do(ctx, http.MethodPost, path, nil, payload, response)
}

// Put performs a PUT request with payload as body, decoding the result into response.
func (c *client) Put(ctx context.Context, path string, payload any, response any) error {
	return c.do(ctx, http.MethodPut, path, nil, payload, response)
}

// Patch performs a PATCH request with payload as body, decoding the result into response.
func (c *client) Patch(ctx context.Context, path string, payload any, response any) error {
	return c.do(ctx, http.MethodPatch, path, nil, payload, response)
}

// Delete performs a DELETE request, decoding the result into response.
func (c *client) Delete(ctx context.Context, path string, response any) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil, response)
}

func (c *client) do(ctx context.Context, method string, path string, query url.Values, payload any, response any) error {
	reqURL, err := url.Parse(c.endpoint)
	if err != nil {
		return api.NewAPIError(fmt.Errorf("invalid IaC API endpoint: %w", err))
	}
	reqURL = reqURL.JoinPath(path)
	reqURL.RawQuery = query.Encode()

	var body io.Reader
	if payload != nil {
		content, err := json.Marshal(payload)
		if err != nil {
			return api.NewAPIError(err)
		}
		body = bytes.NewReader(content)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), body)
	if err != nil {
		return api.NewAPIError(err)
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return api.NewAPIError(err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.NewAPIError(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newResponseError(resp.StatusCode, content)
	}
	if response == nil || len(content) == 0 {
		return nil
	}
	if err := json.Unmarshal(content, response); err != nil {
		return api.NewAPIError(fmt.Errorf("invalid API response: %w", err))
	}
	return nil
}

// resourcePath returns an API path from its parts, escaping each of them.
func resourcePath(parts ...string) string {
	escaped := make([]string, len(parts))
	for i, part := range parts {
		switch part {
		case ".", "..":
			// dots aren't escaped by PathEscape, and dot segments would be
			// removed when joining paths
			escaped[i] = strings.ReplaceAll(part, ".", "%2E")
		default:
			escaped[i] = url.PathEscape(part)
		}
	}
	return strings.Join(escaped, "/")
}
