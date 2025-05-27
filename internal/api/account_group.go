// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// AccountGroup is the data returned by reading account group data.
type AccountGroup struct {
	ID          string
	UUID        string
	Name        string
	Description *string
	Provider    string
	Regions     []string
}

// AccountGroupCreateInput is the input for creating an account group.
type AccountGroupCreateInput struct {
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	Description *string  `json:"description,omitempty"`
	Regions     []string `json:"regions"`
}

func (i AccountGroupCreateInput) GetGraphQLType() string {
	return "AddAccountGroupInput"
}

// AccountGroupUpdateInput is the input for updating an account group.
type AccountGroupUpdateInput struct {
	UUID        string   `json:"uuid"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Regions     []string `json:"regions"`
}

func (i AccountGroupUpdateInput) GetGraphQLType() string {
	return "UpdateAccountGroupInput"
}

type accountGroupAPI struct {
	c *graphql.Client
}

// Read returns data for an account group.
func (a accountGroupAPI) Read(ctx context.Context, uuid string, name string) (AccountGroup, error) {
	var query struct {
		AccountGroup AccountGroup `graphql:"accountGroup(uuid: $uuid, name: $name)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(uuid),
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return query.AccountGroup, NewAPIError(err)
	}

	if query.AccountGroup.ID == "" {
		return query.AccountGroup, NotFound{"Account group not found"}
	}

	return query.AccountGroup, nil
}

// Create creates an account group.
func (a accountGroupAPI) Create(ctx context.Context, i AccountGroupCreateInput) (AccountGroup, error) {
	var mutation struct {
		AddAccountGroup struct {
			Group AccountGroup
		} `graphql:"addAccountGroup(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return mutation.AddAccountGroup.Group, NewAPIError(err)
	}

	return mutation.AddAccountGroup.Group, nil
}

// Update updates an account group.
func (a accountGroupAPI) Update(ctx context.Context, i AccountGroupUpdateInput) (AccountGroup, error) {
	var mutation struct {
		UpdateAccountGroup struct {
			Group AccountGroup
		} `graphql:"updateAccountGroup(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return mutation.UpdateAccountGroup.Group, NewAPIError(err)
	}

	return mutation.UpdateAccountGroup.Group, nil
}

// Delete removes an account group.
func (a accountGroupAPI) Delete(ctx context.Context, uuid string) error {
	var mutation struct {
		RemoveAccountGroup struct {
			Group struct {
				UUID string
			}
		} `graphql:"removeAccountGroup(uuid: $uuid)"`
	}
	variables := map[string]any{"uuid": graphql.String(uuid)}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
}
