// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

const policyCollectionResponse = `{
	"name": "my collection",
	"policies": ["policy-a", "policy-b"],
	"is_default": true,
	"is_mirror": true,
	"created_at": "2026-11-05T16:38:09.334312Z",
	"updated_at": "2026-11-05T16:38:09.334324Z"
}`

func TestPolicyCollectionRead(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, policyCollectionResponse)

	collection, err := policyCollectionAPI{c}.Read(t.Context(), "my collection")

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/base/policy-collections/my%20collection", req.URI)
	assert.Equal(t, &PolicyCollection{
		Name:      "my collection",
		Policies:  []string{"policy-a", "policy-b"},
		IsDefault: true,
		IsMirror:  true,
	}, collection)
}

func TestPolicyCollectionReadNotFound(t *testing.T) {
	c, _ := newTestServer(t, http.StatusNotFound, `{"detail": "Not Found"}`)

	_, err := policyCollectionAPI{c}.Read(t.Context(), "missing")

	assert.Equal(t, api.NotFound{Message: "Not Found"}, err)
}

func TestPolicyCollectionCreate(t *testing.T) {
	tests := []struct {
		name  string
		input PolicyCollectionCreateInput
		body  string
	}{
		{
			name:  "minimal",
			input: PolicyCollectionCreateInput{Name: "my collection"},
			body:  `{"name":"my collection","policies":[],"is_default":false}`,
		},
		{
			name: "all fields",
			input: PolicyCollectionCreateInput{
				Name:      "my collection",
				Policies:  []string{"policy-a", "policy-b"},
				IsDefault: true,
			},
			body: `{"name":"my collection","policies":["policy-a","policy-b"],"is_default":true}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, req := newTestServer(t, http.StatusCreated, policyCollectionResponse)

			collection, err := policyCollectionAPI{c}.Create(t.Context(), tc.input)

			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/base/policy-collections", req.URI)
			assert.JSONEq(t, tc.body, req.Body)
			assert.Equal(t, "my collection", collection.Name)
		})
	}
}

func TestPolicyCollectionUpdate(t *testing.T) {
	tests := []struct {
		name  string
		input PolicyCollectionUpdateInput
		body  string
	}{
		{
			name:  "no policies",
			input: PolicyCollectionUpdateInput{},
			body:  `{"policies":[],"is_default":false}`,
		},
		{
			name: "policies",
			input: PolicyCollectionUpdateInput{
				Policies:  []string{"policy-a"},
				IsDefault: true,
			},
			body: `{"policies":["policy-a"],"is_default":true}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, req := newTestServer(t, http.StatusOK, policyCollectionResponse)

			collection, err := policyCollectionAPI{c}.Update(t.Context(), "my collection", tc.input)

			require.NoError(t, err)
			assert.Equal(t, http.MethodPatch, req.Method)
			assert.Equal(t, "/base/policy-collections/my%20collection", req.URI)
			assert.JSONEq(t, tc.body, req.Body)
			assert.Equal(t, "my collection", collection.Name)
		})
	}
}

func TestPolicyCollectionDelete(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, policyCollectionResponse)

	err := policyCollectionAPI{c}.Delete(t.Context(), "my collection")

	assert.NoError(t, err)
	assert.Equal(t, http.MethodDelete, req.Method)
	assert.Equal(t, "/base/policy-collections/my%20collection", req.URI)
}
