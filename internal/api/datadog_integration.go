// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"
)

// DatadogIntegration is the data for the site-wide Datadog integration.
//
// The API and application keys are absent by design: they are write-only, and
// whether they are stored is reported by Configured rather than by a value.
type DatadogIntegration struct {
	Enabled    bool       `graphql:"enabled"`
	Site       string     `graphql:"site"`
	Configured bool       `graphql:"configured"`
	StatusInfo StatusInfo `graphql:"statusInfo"`
}

// DatadogIntegrationInput is the input for updating the Datadog integration.
//
// Every field is omitted when unset, since the API leaves out fields
// unchanged. That is what lets the site be changed, or the integration
// enabled, without resupplying credentials that are never returned.
type DatadogIntegrationInput struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Site    *string `json:"site,omitempty"`
	APIKey  *string `json:"apiKey,omitempty"`
	AppKey  *string `json:"appKey,omitempty"`
}

func (i DatadogIntegrationInput) GetGraphQLType() string {
	return "UpdateDatadogIntegrationInput"
}

type datadogIntegrationAPI struct {
	c *client
}

// Read returns the site-wide Datadog integration, with defaults when it has
// never been configured.
func (a datadogIntegrationAPI) Read(ctx context.Context) (*DatadogIntegration, error) {
	var query struct {
		DatadogIntegration DatadogIntegration `graphql:"datadogIntegration"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}
	return &query.DatadogIntegration, nil
}

// Update applies changes to the site-wide Datadog integration. Supplying
// either key makes the API check both against Datadog.
func (a datadogIntegrationAPI) Update(ctx context.Context, input DatadogIntegrationInput) (*DatadogIntegration, error) {
	var mutation struct {
		Payload struct {
			DatadogIntegration *DatadogIntegration `graphql:"datadogIntegration"`
			Problems           []problem
		} `graphql:"updateDatadogIntegration(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": input}); err != nil {
		return nil, err
	}
	if err := fromProblems(ctx, mutation.Payload.Problems); err != nil {
		return nil, err
	}
	if mutation.Payload.DatadogIntegration == nil {
		return nil, NotFound{"Datadog integration not found after update"}
	}
	return mutation.Payload.DatadogIntegration, nil
}

// Remove clears the stored credentials and configuration, and tells the
// execution regions to drop their copies.
func (a datadogIntegrationAPI) Remove(ctx context.Context) error {
	var mutation struct {
		Payload struct {
			Problems []problem
		} `graphql:"removeDatadogIntegration"`
	}
	if err := a.c.Mutate(ctx, &mutation, nil); err != nil {
		return err
	}
	return fromProblems(ctx, mutation.Payload.Problems)
}
