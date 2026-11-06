// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

const policySourceResponse = `{
	"name": "my source",
	"url": "https://github.com/example/repo",
	"branch": "main",
	"paths": ["policies"],
	"extensions": [".yaml"],
	"auto_scan": true,
	"hook_id": "abc",
	"hook_path": "/policy-sources/my%20source/webhooks/abc",
	"auth": {"username": "user", "token": "**********"}
}`

func ptr[T any](v T) *T {
	return &v
}

func TestPolicySourceRead(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, policySourceResponse)

	source, err := policySourceAPI{c}.Read(t.Context(), "my source")

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/base/policy-sources/my%20source", req.URI)
	assert.Equal(t, &PolicySource{
		Name:       "my source",
		URL:        "https://github.com/example/repo",
		Branch:     ptr("main"),
		Paths:      []string{"policies"},
		Extensions: []string{".yaml"},
		AutoScan:   true,
		HookPath:   "/policy-sources/my%20source/webhooks/abc",
		Auth:       &PolicySourceAuth{Username: "user", Token: ptr("**********")},
		WebhookURL: c.endpoint + "policy-sources/my%20source/webhooks/abc",
	}, source)
}

func TestPolicySourceReadNotFound(t *testing.T) {
	c, _ := newTestServer(t, http.StatusNotFound, `{"detail": "Not Found"}`)

	_, err := policySourceAPI{c}.Read(t.Context(), "missing")

	assert.Equal(t, api.NotFound{Message: "Not Found"}, err)
}

func TestPolicySourceCreate(t *testing.T) {
	tests := []struct {
		name  string
		input PolicySourceCreateInput
		body  string
	}{
		{
			name: "minimal",
			input: PolicySourceCreateInput{
				Name: "my source",
				URL:  "https://github.com/example/repo",
			},
			body: `{"name":"my source","url":"https://github.com/example/repo","auto_scan":false}`,
		},
		{
			name: "token auth",
			input: PolicySourceCreateInput{
				Name:       "my source",
				URL:        "https://github.com/example/repo",
				Branch:     ptr("main"),
				Paths:      []string{"policies"},
				Extensions: []string{".yaml"},
				AutoScan:   true,
				Auth:       &PolicySourceAuth{Username: "user", Token: ptr("secret")},
			},
			body: `{"name":"my source","url":"https://github.com/example/repo","branch":"main","paths":["policies"],"extensions":[".yaml"],"auto_scan":true,"auth":{"username":"user","token":"secret"}}`,
		},
		{
			name: "ssh auth",
			input: PolicySourceCreateInput{
				Name: "my source",
				URL:  "ssh://git@github.com/example/repo",
				Auth: &PolicySourceAuth{Username: "git", SSHKey: ptr("key"), SSHPassphrase: ptr("pass")},
			},
			body: `{"name":"my source","url":"ssh://git@github.com/example/repo","auto_scan":false,"auth":{"username":"git","ssh_key":"key","ssh_passphrase":"pass"}}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, req := newTestServer(t, http.StatusCreated, policySourceResponse)

			source, err := policySourceAPI{c}.Create(t.Context(), tc.input)

			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/base/policy-sources", req.URI)
			assert.JSONEq(t, tc.body, req.Body)
			assert.Equal(t, "my source", source.Name)
			assert.Equal(t, c.endpoint+"policy-sources/my%20source/webhooks/abc", source.WebhookURL)
		})
	}
}

func TestPolicySourceUpdate(t *testing.T) {
	tests := []struct {
		name  string
		input PolicySourceUpdateInput
		body  string
	}{
		{
			name: "auth unchanged",
			input: PolicySourceUpdateInput{
				URL:        "https://github.com/example/repo",
				Extensions: []string{".yaml"},
			},
			body: `{"url":"https://github.com/example/repo","branch":null,"paths":null,"extensions":[".yaml"],"auto_scan":false}`,
		},
		{
			name: "auth removed",
			input: PolicySourceUpdateInput{
				URL:        "https://github.com/example/repo",
				Extensions: []string{".yaml"},
				Auth:       PolicySourceAuthUpdate{Set: true},
			},
			body: `{"url":"https://github.com/example/repo","branch":null,"paths":null,"extensions":[".yaml"],"auto_scan":false,"auth":null}`,
		},
		{
			name: "auth set",
			input: PolicySourceUpdateInput{
				URL:        "https://github.com/example/repo",
				Branch:     ptr("main"),
				Paths:      []string{"policies"},
				Extensions: []string{".yaml"},
				AutoScan:   true,
				Auth: PolicySourceAuthUpdate{
					Set:  true,
					Auth: &PolicySourceAuth{Username: "user", Token: ptr("secret")},
				},
			},
			body: `{"url":"https://github.com/example/repo","branch":"main","paths":["policies"],"extensions":[".yaml"],"auto_scan":true,"auth":{"username":"user","token":"secret"}}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// the update response doesn't include auth, the source is read
			// back to get it
			var requests []capturedRequest
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				requests = append(requests, capturedRequest{
					Method: r.Method,
					URI:    r.URL.RequestURI(),
					Body:   string(body),
				})
				if r.Method == http.MethodPatch {
					_, _ = w.Write([]byte(`{"name": "my source", "auth": null}`))
				} else {
					_, _ = w.Write([]byte(policySourceResponse))
				}
			}))
			t.Cleanup(server.Close)
			c := newClient(ClientConfig{Endpoint: server.URL + "/base/", APIKey: "secret"})

			source, err := policySourceAPI{c}.Update(t.Context(), "my source", tc.input)

			require.NoError(t, err)
			require.Len(t, requests, 2)
			assert.Equal(t, http.MethodPatch, requests[0].Method)
			assert.Equal(t, "/base/policy-sources/my%20source", requests[0].URI)
			assert.JSONEq(t, tc.body, requests[0].Body)
			assert.Equal(t, http.MethodGet, requests[1].Method)
			assert.Equal(t, "/base/policy-sources/my%20source", requests[1].URI)
			assert.Equal(t, "my source", source.Name)
			assert.Equal(t, &PolicySourceAuth{Username: "user", Token: ptr("**********")}, source.Auth)
		})
	}
}

func TestPolicySourceDelete(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, policySourceResponse)

	err := policySourceAPI{c}.Delete(t.Context(), "my source")

	assert.NoError(t, err)
	assert.Equal(t, http.MethodDelete, req.Method)
	assert.Equal(t, "/base/policy-sources/my%20source", req.URI)
}
