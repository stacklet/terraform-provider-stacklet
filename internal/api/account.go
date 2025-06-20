// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Account is the data returned by reading account data.
type Account struct {
	ID              string
	Key             string
	Name            string
	ShortName       *string
	Description     *string
	Provider        CloudProvider
	Path            *string
	Email           *string
	Active          bool
	SecurityContext *string
	Variables       *string
}

// AccountCreateInput is the input for creating an account.
type AccountCreateInput struct {
	Name            string        `json:"name"`
	Key             string        `json:"key"`
	Provider        CloudProvider `json:"provider"`
	ShortName       *string       `json:"shortName,omitempty"`
	Description     *string       `json:"description,omitempty"`
	Email           *string       `json:"email,omitempty"`
	SecurityContext *string       `json:"securityContext,omitempty"`
	Variables       *string       `json:"variables,omitempty"`
}

func (i AccountCreateInput) GetGraphQLType() string {
	return "AccountInput"
}

// AccountUpdateInput is the input for updating an account.
type AccountUpdateInput struct {
	Key             string        `json:"key"`
	Provider        CloudProvider `json:"provider"`
	Name            *string       `json:"name"`
	ShortName       *string       `json:"shortName"`
	Description     *string       `json:"description"`
	Email           *string       `json:"email"`
	SecurityContext *string       `json:"securityContext,omitempty"`
	Variables       *string       `json:"variables"`
}

func (i AccountUpdateInput) GetGraphQLType() string {
	return "UpdateAccountInput"
}

type accountAPI struct {
	c *graphql.Client
}

// Read returns data for an account.
func (a accountAPI) Read(ctx context.Context, cloudProvider string, key string) (*Account, error) {
	var query struct {
		Account Account `graphql:"account(provider: $provider, key: $key)"`
	}
	variables := map[string]any{
		"provider": CloudProvider(cloudProvider),
		"key":      graphql.String(key),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if query.Account.ID == "" || !query.Account.Active {
		return nil, NotFound{"Account not found"}
	}

	return &query.Account, nil
}

// Create creates an account.
func (a accountAPI) Create(ctx context.Context, i AccountCreateInput) (*Account, error) {
	var mutation struct {
		Payload struct {
			Account Account
		} `graphql:"addAccount(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Account, nil
}

// Update updates an account.
func (a accountAPI) Update(ctx context.Context, i AccountUpdateInput) (*Account, error) {
	var mutation struct {
		Payload struct {
			Account Account
		} `graphql:"updateAccount(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Account, nil
}

// Delete removes an account.
func (a accountAPI) Delete(ctx context.Context, cloudProvider string, key string) error {
	var mutation struct {
		Payload struct {
			Account struct {
				Key string
			}
		} `graphql:"removeAccount(provider: $provider, key: $key)"`
	}
	variables := map[string]any{
		"provider": CloudProvider(cloudProvider),
		"key":      graphql.String(key),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
}
