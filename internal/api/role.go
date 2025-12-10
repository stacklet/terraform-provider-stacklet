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

// List returns all roles.
func (r roleAPI) List(ctx context.Context) ([]Role, error) {
	cursor := ""
	var query struct {
		Roles struct {
			Edges []struct {
				Node Role
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"roles(first: 100, after: $cursor)"`
	}

	roles := make([]Role, 0)

	// Paginate through all results
	for {
		variables := map[string]any{
			"cursor": graphql.String(cursor),
		}

		if err := r.c.Query(ctx, &query, variables); err != nil {
			return nil, NewAPIError(err)
		}

		for _, edge := range query.Roles.Edges {
			roles = append(roles, edge.Node)
		}

		// Check if there are more pages
		if !query.Roles.PageInfo.HasNextPage {
			break
		}
		cursor = query.Roles.PageInfo.EndCursor
	}

	return roles, nil
}
