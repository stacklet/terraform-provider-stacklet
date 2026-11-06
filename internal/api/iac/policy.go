// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
)

// PolicyResources is the list of resource types a policy applies to.
//
// The API returns either a single resource type or a list of them.
type PolicyResources []string

func (r *PolicyResources) UnmarshalJSON(data []byte) error {
	var resources []string
	if err := json.Unmarshal(data, &resources); err == nil {
		*r = resources
		return nil
	}
	var resource string
	if err := json.Unmarshal(data, &resource); err != nil {
		return err
	}
	*r = PolicyResources{resource}
	return nil
}

// Policy is an IaC policy, found by scanning a policy source.
type Policy struct {
	Name       string          `json:"name"`
	SourceName string          `json:"source_name"`
	Resources  PolicyResources `json:"resource"`
	Severity   string          `json:"severity"`
}

type policyAPI struct {
	c *client
}

// ListBySource returns policies from a policy source, sorted by name.
func (a policyAPI) ListBySource(ctx context.Context, source string) ([]Policy, error) {
	return a.list(ctx, resourcePath("policy-sources", source, "policies"))
}

// ListByCollection returns policies in a policy collection, sorted by name.
func (a policyAPI) ListByCollection(ctx context.Context, collection string) ([]Policy, error) {
	return a.list(ctx, resourcePath("policy-collections", collection, "policies"))
}

func (a policyAPI) list(ctx context.Context, path string) ([]Policy, error) {
	policies := []Policy{}
	if err := a.c.Get(ctx, path, nil, &policies); err != nil {
		return nil, err
	}
	// the API doesn't guarantee ordering
	slices.SortFunc(policies, func(a, b Policy) int {
		return strings.Compare(a.Name, b.Name)
	})
	return policies, nil
}
