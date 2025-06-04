// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// BindingExecutionConfig holds the execution configuration for a binding.
type BindingExecutionConfig struct {
	DryRun    *BindingExecutionConfigDryRun `json:"dryRun"`
	Variables *string                       `json:"variables"`
}

// DryRunDefault returns the dry run value.
func (c BindingExecutionConfig) DryRunDefault() bool {
	if c.DryRun == nil {
		return false
	}
	return c.DryRun.Default
}

// BindingExecutionConfigDryRun holds the dry run confiuration for a binding execution configuration.
type BindingExecutionConfigDryRun struct {
	Default bool `json:"default"`
}

// BindingExecutionConfigUpdateInput is the data to update a binding execution configuration.
type BindingExecutionConfigUpdateInput struct {
	BindingUUID     string                 `json:"uuid"`
	ExecutionConfig BindingExecutionConfig `json:"executionConfig"`
}

func (i BindingExecutionConfigUpdateInput) GetGraphQLType() string {
	return "UpdateBindingInput"
}

type bindingExecutionConfigAPI struct {
	c *graphql.Client
}

// Read returns data for a binding execution configuration.
func (a bindingExecutionConfigAPI) Read(ctx context.Context, uuid string) (*BindingExecutionConfig, error) {
	var query struct {
		Binding struct {
			ID              string
			ExecutionConfig BindingExecutionConfig
		} `graphql:"binding(uuid: $uuid)"`
	}
	variables := map[string]any{"uuid": graphql.String(uuid)}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}
	if query.Binding.ID == "" {
		return nil, NotFound{"Binding not found"}
	}

	return &query.Binding.ExecutionConfig, nil
}

// Upsert updates a binding execution configuration.
func (a bindingExecutionConfigAPI) Upsert(ctx context.Context, i BindingExecutionConfigUpdateInput) (*BindingExecutionConfig, error) {
	var mutation struct {
		Payload struct {
			Binding struct {
				ExecutionConfig BindingExecutionConfig
			}
		} `graphql:"updateBinding(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Binding.ExecutionConfig, nil
}
