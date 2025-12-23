// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// User is the data returned by reading user data.
type User struct {
	ID                      string
	Active                  bool
	AllRoles                []string
	AssignedRoles           []string
	DisplayName             *string
	Email                   *string
	Groups                  []string
	ImplicitRoles           []string
	InheritedRoles          []string
	Key                     int
	LastLogin               *string
	Name                    *string
	RoleAssignmentPrincipal string
	Roles                   []string
	SSOUser                 bool
	Username                *string
}

type userAPI struct {
	c *graphql.Client
}

// Read returns data for a user by username.
// Note: The Stacklet API user filter requires the field name "username" with no operator.
func (u userAPI) Read(ctx context.Context, username string) (*User, error) {
	var query struct {
		Users struct {
			Edges []struct {
				Node User
			}
		} `graphql:"users(filterElement: $filterElement)"`
	}
	// Use "username" as field name and omit operator
	variables := map[string]any{
		"filterElement": newSimpleFilter("username", username),
	}
	if err := u.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(query.Users.Edges) == 0 {
		return nil, NotFound{"User not found"}
	}

	return &query.Users.Edges[0].Node, nil
}
