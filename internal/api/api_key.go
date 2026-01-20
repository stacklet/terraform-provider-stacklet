// Copyright (c) 2026 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// APIKey is the data returned by reading API key details.
type APIKey struct {
	ID          string
	Identity    string
	Description string
	ExpiresAt   *string
	RevokedAt   *string
}

// AddAPIKeyInput is the input for creating an API key.
type AddAPIKeyInput struct {
	Description string  `json:"description"`
	ExpiresAt   *string `json:"expiresAt,omitempty"`
}

// UpdateAPIKeyInput is the input for updating an API key.
type UpdateAPIKeyInput struct {
	Identity    string `json:"identity"`
	Description string `json:"description"`
}

// RevokeyAPIKeysInput is the input for revoking API keys.
type RevokeAPIKeysInput struct {
	Identities []string `json:"identities"`
}

type apiKeyAPI struct {
	c *graphql.Client
}

// Read returns data for an API Key.
func (a apiKeyAPI) Read(ctx context.Context, identity string) (*APIKey, error) {
	var query struct {
		Conn struct {
			Edges []struct {
				Node APIKey
			}
		} `graphql:"apiKeys(filterElement: $filterElement)"`
	}

	variables := map[string]any{
		"filterElement": newSimpleFilter("identity", identity),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(query.Conn.Edges) == 0 {
		return nil, NotFound{"API Key not found"}
	}

	return &query.Conn.Edges[0].Node, nil
}

// Create creates an API key.
func (a apiKeyAPI) Create(ctx context.Context, i AddAPIKeyInput) (*APIKey, *string, error) {
	var mutation struct {
		Payload struct {
			Key      APIKey
			Secret   string
			Problems []Problem
		} `graphql:"addAPIKey(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, nil, NewAPIError(err)
	}
	if err := fromProblems(ctx, mutation.Payload.Problems); err != nil {
		return nil, nil, err
	}

	return &mutation.Payload.Key, &mutation.Payload.Secret, nil
}

// Update updates an API Key.
func (a apiKeyAPI) Update(ctx context.Context, i UpdateAPIKeyInput) (*APIKey, error) {
	var mutation struct {
		Payload struct {
			Key      APIKey
			Problems []Problem
		} `graphql:"updateAPIKey(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}
	if err := fromProblems(ctx, mutation.Payload.Problems); err != nil {
		return nil, err
	}

	return &mutation.Payload.Key, nil
}

// Revoke revokes an API key.
func (a apiKeyAPI) Revoke(ctx context.Context, identity string) error {
	var mutation struct {
		Payload struct {
			Problems []Problem
		} `graphql:"revokeAPIKeys(input: $input)"`
	}
	input := map[string]any{"input": RevokeAPIKeysInput{Identities: []string{identity}}}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return NewAPIError(err)
	}
	return fromProblems(ctx, mutation.Payload.Problems)
}
