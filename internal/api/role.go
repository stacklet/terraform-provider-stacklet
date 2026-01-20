// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Role is the data returned by reading role data.
type Role struct {
	ID          string
	Name        string
	Permissions []string
	System      bool
}

type roleAPI struct {
	c *graphql.Client
}

// Read returns data for a role by name.
func (r roleAPI) Read(ctx context.Context, name string) (*Role, error) {
	var query struct {
		Roles struct {
			Edges []struct {
				Node Role
			}
		} `graphql:"roles(filterElement: $filterElement)"`
	}
	variables := map[string]any{
		"filterElement": newExactMatchFilter("name", name),
	}
	if err := r.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(query.Roles.Edges) == 0 {
		return nil, NotFound{"Role not found"}
	}

	return &query.Roles.Edges[0].Node, nil
}
