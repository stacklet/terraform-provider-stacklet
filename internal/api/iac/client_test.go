// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

type item struct {
	Name string `json:"name"`
}

type capturedRequest struct {
	Method      string
	URI         string
	Body        string
	ContentType string
	Auth        string
}

func newTestServer(t *testing.T, status int, response string) (*client, *capturedRequest) {
	captured := &capturedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*captured = capturedRequest{
			Method:      r.Method,
			URI:         r.URL.RequestURI(),
			Body:        string(body),
			ContentType: r.Header.Get("Content-Type"),
			Auth:        r.Header.Get("Authorization"),
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	c := newClient(ClientConfig{Endpoint: server.URL + "/base/", APIKey: "secret", Version: "1.0.0"})
	return c, captured
}

func TestClientGet(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, `{"name": "foo"}`)

	var response item
	err := c.Get(t.Context(), resourcePath("projects", "a/b c"), url.Values{"page": {"2"}}, &response)

	assert.NoError(t, err)
	assert.Equal(t, item{Name: "foo"}, response)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/base/projects/a%2Fb%20c?page=2", req.URI)
	assert.Equal(t, "Bearer secret", req.Auth)
	assert.Equal(t, "", req.ContentType)
}

func TestClientDotSegments(t *testing.T) {
	for _, name := range []string{".", ".."} {
		t.Run(name, func(t *testing.T) {
			c, req := newTestServer(t, http.StatusNotFound, `{"detail": "Not found"}`)

			err := c.Get(t.Context(), resourcePath("projects", name), nil, nil)

			assert.Equal(t, api.NotFound{Message: "Not found"}, err)
			assert.Equal(t, "/base/projects/"+strings.ReplaceAll(name, ".", "%2E"), req.URI)
		})
	}
}

func TestClientWithBody(t *testing.T) {
	tests := []struct {
		method string
		call   func(c *client, response *item) error
	}{
		{
			method: http.MethodPost,
			call: func(c *client, response *item) error {
				return c.Post(t.Context(), "projects", item{Name: "foo"}, response)
			},
		},
		{
			method: http.MethodPut,
			call: func(c *client, response *item) error {
				return c.Put(t.Context(), "projects", item{Name: "foo"}, response)
			},
		},
		{
			method: http.MethodPatch,
			call: func(c *client, response *item) error {
				return c.Patch(t.Context(), "projects", item{Name: "foo"}, response)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			c, req := newTestServer(t, http.StatusCreated, `{"name": "foo"}`)

			var response item
			err := tc.call(c, &response)

			assert.NoError(t, err)
			assert.Equal(t, item{Name: "foo"}, response)
			assert.Equal(t, tc.method, req.Method)
			assert.Equal(t, "/base/projects", req.URI)
			assert.Equal(t, `{"name":"foo"}`, req.Body)
			assert.Equal(t, "application/json", req.ContentType)
		})
	}
}

func TestClientDeleteNoContent(t *testing.T) {
	c, req := newTestServer(t, http.StatusNoContent, "")

	err := c.Delete(t.Context(), "projects/foo", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.MethodDelete, req.Method)
	assert.Equal(t, "/base/projects/foo", req.URI)
}

func TestClientInvalidResponse(t *testing.T) {
	c, _ := newTestServer(t, http.StatusOK, "not json")

	var response item
	err := c.Get(t.Context(), "projects/foo", nil, &response)

	assert.ErrorContains(t, err, "invalid API response")
}

func TestClientErrors(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		response string
		err      error
	}{
		{
			name:     "not found",
			status:   http.StatusNotFound,
			response: `{"detail": "Project not found"}`,
			err:      api.NotFound{Message: "Project not found"},
		},
		{
			name:     "unauthorized",
			status:   http.StatusUnauthorized,
			response: `{"detail": "Invalid Authorization token"}`,
			err:      api.APIError{Kind: "Authorization Error", Detail: "Invalid Authorization token"},
		},
		{
			name:     "gateway message",
			status:   http.StatusForbidden,
			response: `{"message": "Forbidden"}`,
			err:      api.APIError{Kind: "Authorization Error", Detail: "Forbidden"},
		},
		{
			name:     "validation",
			status:   http.StatusUnprocessableEntity,
			response: `{"detail": [{"loc": ["body", "name"], "msg": "Field required"}, {"loc": ["body", "groups", 0], "msg": "Invalid"}]}`,
			err:      api.APIError{Kind: "API Error", Detail: "body.name: Field required\nbody.groups.0: Invalid"},
		},
		{
			name:     "not json",
			status:   http.StatusInternalServerError,
			response: "Internal server error",
			err:      api.APIError{Kind: "API Error", Detail: "Internal server error"},
		},
		{
			name:     "empty",
			status:   http.StatusBadGateway,
			response: "",
			err:      api.APIError{Kind: "API Error", Detail: "Bad Gateway"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := newTestServer(t, tc.status, tc.response)

			err := c.Get(t.Context(), "projects/foo", nil, nil)

			assert.Equal(t, tc.err, err)
		})
	}
}
