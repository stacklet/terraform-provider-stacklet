// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

const policiesResponse = `[
	{
		"name": "policy-b",
		"source_name": "my source",
		"raw_policy": {"name": "policy-b", "resource": ["terraform.aws_s3_bucket", "terraform.aws_s3_object"]},
		"resource": ["terraform.aws_s3_bucket", "terraform.aws_s3_object"],
		"severity": "high",
		"created_at": "2026-11-05T16:38:09.334312Z",
		"updated_at": "2026-11-05T16:38:09.334324Z"
	},
	{
		"name": "policy-a",
		"source_name": "my source",
		"raw_policy": {"name": "policy-a", "resource": "terraform.aws_instance"},
		"resource": "terraform.aws_instance",
		"severity": "low",
		"created_at": "2026-11-05T16:38:09.334312Z",
		"updated_at": "2026-11-05T16:38:09.334324Z"
	}
]`

var expectedPolicies = []Policy{
	{
		Name:       "policy-a",
		SourceName: "my source",
		Resources:  PolicyResources{"terraform.aws_instance"},
		Severity:   "low",
	},
	{
		Name:       "policy-b",
		SourceName: "my source",
		Resources:  PolicyResources{"terraform.aws_s3_bucket", "terraform.aws_s3_object"},
		Severity:   "high",
	},
}

func TestPolicyListBySource(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, policiesResponse)

	policies, err := policyAPI{c}.ListBySource(t.Context(), "my source")

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/base/policy-sources/my%20source/policies", req.URI)
	assert.Equal(t, expectedPolicies, policies)
}

func TestPolicyListByCollection(t *testing.T) {
	c, req := newTestServer(t, http.StatusOK, policiesResponse)

	policies, err := policyAPI{c}.ListByCollection(t.Context(), "my collection")

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/base/policy-collections/my%20collection/policies", req.URI)
	assert.Equal(t, expectedPolicies, policies)
}

func TestPolicyListEmpty(t *testing.T) {
	c, _ := newTestServer(t, http.StatusOK, `[]`)

	policies, err := policyAPI{c}.ListBySource(t.Context(), "my source")

	require.NoError(t, err)
	assert.Equal(t, []Policy{}, policies)
}

func TestPolicyListNotFound(t *testing.T) {
	c, _ := newTestServer(t, http.StatusNotFound, `{"detail": "Source not found"}`)

	_, err := policyAPI{c}.ListBySource(t.Context(), "missing")

	assert.Equal(t, api.NotFound{Message: "Source not found"}, err)
}

func TestPolicyResourcesUnmarshalInvalid(t *testing.T) {
	var r PolicyResources
	assert.Error(t, r.UnmarshalJSON([]byte(`123`)))
}
