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
		Role Role `graphql:"role(name: $name)"`
	}
	variables := map[string]any{
		"name": graphql.String(name),
	}
	if err := r.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if query.Role.ID == "" {
		return nil, NotFound{"Role not found"}
	}

	return &query.Role, nil
}

// List returns all roles.
func (r roleAPI) List(ctx context.Context) ([]Role, error) {
	var query struct {
		Roles struct {
			Edges []struct {
				Node Role
			}
		}
	}
	if err := r.c.Query(ctx, &query, nil); err != nil {
		return nil, NewAPIError(err)
	}

	roles := make([]Role, 0, len(query.Roles.Edges))
	for _, edge := range query.Roles.Edges {
		roles = append(roles, edge.Node)
	}

	return roles, nil
}
