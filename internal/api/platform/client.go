// Copyright Stacklet, Inc. 2025, 2026

package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// ClientConfig is the configuration for the API client.
type ClientConfig struct {
	Endpoint string
	APIKey   string
	Version  string
	PageSize int
}

// client is the wrapper for the GraphQL client.
type client struct {
	c        *graphql.Client
	pageSize int
}

// Query makes a GraphQL query call.
func (c *client) Query(ctx context.Context, q any, variables map[string]any) error {
	err := c.c.Query(ctx, q, variables)
	if err != nil {
		return api.NewAPIError(err)
	}
	return nil
}

// Mutate makes a GraphQL mutation call.
func (c *client) Mutate(ctx context.Context, m any, variables map[string]any) error {
	err := c.c.Mutate(ctx, m, variables)
	if err != nil {
		return api.NewAPIError(err)
	}
	return nil
}

// newClient returns a configured graphql Client.
func newClient(config ClientConfig) *client {
	httpClient := &http.Client{
		Transport: &errorTransport{
			Base: api.NewTransport(config.APIKey, config.Version),
		},
	}
	return &client{
		c:        graphql.NewClient(config.Endpoint, httpClient),
		pageSize: config.PageSize,
	}
}

// errorTransport is an http.Transport that converts 400 responses with a
// GraphQL-style error payload into Go errors, so callers see the actual
// error messages rather than a generic "400 Bad Request".
type errorTransport struct {
	Base http.RoundTripper
}

func (t *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)
	if err != nil || resp == nil || resp.StatusCode != http.StatusBadRequest {
		return resp, err
	}

	content, err := api.DecodeResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if json.Unmarshal([]byte(content), &payload) == nil && len(payload.Errors) > 0 {
		msgs := make([]string, len(payload.Errors))
		for i, e := range payload.Errors {
			msgs[i] = e.Message
		}
		return nil, api.APIError{Kind: "API Error", Detail: strings.Join(msgs, "\n")}
	}

	return resp, nil
}
