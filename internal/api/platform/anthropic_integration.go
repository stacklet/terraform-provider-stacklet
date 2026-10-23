// Copyright Stacklet, Inc. 2025, 2026

package platform

import (
	"context"

	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// AnthropicIntegration is the data for an Anthropic organization integration.
//
// The admin key is absent by design: it is write-only, and whether one is
// stored is reported by Configured rather than by a value.
type AnthropicIntegration struct {
	ID         graphql.ID `graphql:"id"`
	OrgID      string     `graphql:"orgId"`
	Configured bool       `graphql:"configured"`
}

// AnthropicIntegrationInput is the input for creating or updating an
// Anthropic integration.
//
// The admin key is omitted when unset: the API then keeps the stored one,
// which is how an existing integration is saved without resupplying a key
// that is never returned.
type AnthropicIntegrationInput struct {
	OrgID    string  `json:"orgId"`
	AdminKey *string `json:"adminKey,omitempty"`
}

func (i AnthropicIntegrationInput) GetGraphQLType() string {
	return "UpsertAnthropicIntegrationInput"
}

type anthropicIntegrationRemoveInput struct {
	OrgID string `json:"orgId"`
}

func (i anthropicIntegrationRemoveInput) GetGraphQLType() string {
	return "RemoveAnthropicIntegrationInput"
}

type anthropicIntegrationAPI struct {
	c *client
}

// Read returns the integration for an Anthropic organization.
func (a anthropicIntegrationAPI) Read(ctx context.Context, orgID string) (*AnthropicIntegration, error) {
	var query struct {
		Payload struct {
			AnthropicIntegration *AnthropicIntegration `graphql:"anthropicIntegration"`
			Problems             []problem
		} `graphql:"anthropicIntegration(orgId: $orgId)"`
	}
	if err := a.c.Query(ctx, &query, map[string]any{"orgId": graphql.String(orgID)}); err != nil {
		return nil, err
	}
	if err := fromProblems(ctx, query.Payload.Problems); err != nil {
		return nil, err
	}
	if query.Payload.AnthropicIntegration == nil {
		return nil, api.NotFound{Message: "Anthropic integration not found"}
	}
	return query.Payload.AnthropicIntegration, nil
}

// Upsert creates or updates the integration for an Anthropic organization.
// A supplied admin key is checked against Anthropic before it is stored, and
// a key that fails the check is reported as a problem.
func (a anthropicIntegrationAPI) Upsert(ctx context.Context, input AnthropicIntegrationInput) (*AnthropicIntegration, error) {
	var mutation struct {
		Payload struct {
			AnthropicIntegration *AnthropicIntegration `graphql:"anthropicIntegration"`
			Problems             []problem
		} `graphql:"upsertAnthropicIntegration(input: $input)"`
	}
	ctx = api.MaskSecrets(ctx, input.AdminKey)
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": input}); err != nil {
		return nil, err
	}
	if err := fromProblems(ctx, mutation.Payload.Problems); err != nil {
		return nil, err
	}
	if mutation.Payload.AnthropicIntegration == nil {
		return nil, api.NotFound{Message: "Anthropic integration not found after upsert"}
	}
	return mutation.Payload.AnthropicIntegration, nil
}

// Remove deletes the integration for an Anthropic organization, along with
// its stored admin key.
func (a anthropicIntegrationAPI) Remove(ctx context.Context, orgID string) error {
	var mutation struct {
		Payload struct {
			Problems []problem
		} `graphql:"removeAnthropicIntegration(input: $input)"`
	}
	input := anthropicIntegrationRemoveInput{OrgID: orgID}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": input}); err != nil {
		return err
	}
	return fromProblems(ctx, mutation.Payload.Problems)
}
