// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Policy is the data returned by reading policy data.
type Policy struct {
	ID              graphql.ID
	UUID            string
	Name            string
	Description     *string
	Provider        string
	Version         int
	Category        []string
	Mode            string
	ResourceType    string
	Path            string
	Source          string
	SourceYAML      string `graphql:"sourceYAML"`
	System          bool
	UnqualifiedName string
}

type policyAPI struct {
	c *client
}

// Read returns data for a policy.
func (a policyAPI) Read(ctx context.Context, uuid string, name string, version int) (*Policy, error) {
	var query struct {
		Policy Policy `graphql:"policy(uuid: $uuid, name: $name, version: $version)"`
	}
	variables := map[string]any{
		"uuid":    graphql.String(uuid),
		"name":    graphql.String(name),
		"version": graphql.Int(version),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	if query.Policy.ID == "" {
		return nil, NotFound{"Policy not found"}
	}

	return &query.Policy, nil
}
