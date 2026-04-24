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
	DisplayName             *string
	Email                   *string
	Name                    *string
	Key                     int64
	RoleAssignmentPrincipal string
	SSOUser                 bool
	Username                *string
}

// UserCreateInput is the input for creating a user.
type UserCreateInput struct {
	Name        string   `json:"name"`
	DisplayName *string  `json:"displayName,omitempty"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	SSOUser     bool     `json:"ssoUser"`
}

func (i UserCreateInput) GetGraphQLType() string {
	return "AddUserInput"
}

// UserUpdateInput is the input for updating a user.
type UserUpdateInput struct {
	Name        *string `json:"name"`
	DisplayName *string `json:"displayName"`
	Key         int64   `json:"key"`
	Email       string  `json:"email"`
}

func (i UserUpdateInput) GetGraphQLType() string {
	return "UpdateUserInput"
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

// Create creates a user.
func (a userAPI) Create(ctx context.Context, i UserCreateInput) (*User, error) {
	var mutation struct {
		Payload struct {
			User User
		} `graphql:"addUser(input: $input)"`
	}

	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.User, nil
}

// Update updates a user.
func (a userAPI) Update(ctx context.Context, i UserUpdateInput) (*User, error) {
	var mutation struct {
		Payload struct {
			User User
		} `graphql:"updateUser(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.User, nil
}

// Delete removes a user.
func (a userAPI) Delete(ctx context.Context, key int64) error {
	var mutation struct {
		Payload struct {
			Removed []struct {
				ID string
			}
		} `graphql:"removeUser(key: $key)"`
	}
	variables := map[string]any{"key": graphql.Int(key)}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
}
