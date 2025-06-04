// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Binding is the data returned by reading binding details.
type Binding struct {
	ID           string
	UUID         string
	Name         string
	Description  *string
	AutoDeploy   bool
	Schedule     *string
	AccountGroup struct {
		UUID string
	}
	PolicyCollection struct {
		UUID string
	}
	ExecutionConfig BindingExecutionConfig
	System          bool
}

// ExecutionConfig holds the execution configuration for a binding.
type BindingExecutionConfig struct {
	DryRun          *BindingExecutionConfigDryRun          `json:"dryRun"`
	SecurityContext *BindingExecutionConfigSecurityContext `json:"securityContext"`
	Variables       *string                                `json:"variables"`
}

// DryRunDefault returns the dry run default value.
func (c BindingExecutionConfig) DryRunDefault() bool {
	if c.DryRun == nil {
		return false
	}
	return c.DryRun.Default
}

// SecurityContextDefault returns the security context default, or nil if not set.
func (c BindingExecutionConfig) SecurityContextDefault() *string {
	if c.SecurityContext == nil {
		return nil
	}
	return &c.SecurityContext.Default
}

// BindingExecutionConfigDryRun holds the dry run confiuration for a binding execution config.
type BindingExecutionConfigDryRun struct {
	Default bool `json:"default"`
}

// BindingExecutionConfigSecurityCotnext holds the security context configuration for a binding execution config.
type BindingExecutionConfigSecurityContext struct {
	Default string `json:"default"`
}

// BindingCreateInput is the input for creating a binding.
type BindingCreateInput struct {
	Name                 string                 `json:"name"`
	Description          *string                `json:"description,omitempty"`
	AutoDeploy           bool                   `json:"autoDeploy"`
	Schedule             *string                `json:"schedule,omitempty"`
	ExecutionConfig      BindingExecutionConfig `json:"executionConfig"`
	AccountGroupUUID     string                 `json:"accountGroupUUID"`
	PolicyCollectionUUID string                 `json:"policyCollectionUUID"`
	Deploy               bool                   `json:"deploy"`
}

func (i BindingCreateInput) GetGraphQLType() string {
	return "AddBindingInput"
}

type BindingUpdateInput struct {
	UUID            string                 `json:"uuid"`
	Name            string                 `json:"name"`
	Description     *string                `json:"description"`
	AutoDeploy      bool                   `json:"autoDeploy"`
	Schedule        *string                `json:"schedule"`
	ExecutionConfig BindingExecutionConfig `json:"executionConfig"`
}

func (i BindingUpdateInput) GetGraphQLType() string {
	return "UpdateBindingInput"
}

type bindingAPI struct {
	c *graphql.Client
}

// Read returns data for a binding.
func (a bindingAPI) Read(ctx context.Context, uuid string, name string) (*Binding, error) {
	var query struct {
		Binding Binding `graphql:"binding(uuid: $uuid, name: $name)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}
	if query.Binding.ID == "" {
		return nil, NotFound{"Binding not found"}
	}

	return &query.Binding, nil
}

// Create creates a binding.
func (a bindingAPI) Create(ctx context.Context, i BindingCreateInput) (*Binding, error) {
	var mutation struct {
		Payload struct {
			Binding Binding
		} `graphql:"addBinding(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Binding, nil
}

// Update updates a binding.
func (a bindingAPI) Update(ctx context.Context, i BindingUpdateInput) (*Binding, error) {
	var mutation struct {
		Payload struct {
			Binding Binding
		} `graphql:"updateBinding(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Binding, nil
}

// Delete removes a binding.
func (a bindingAPI) Delete(ctx context.Context, uuid string) error {
	var mutation struct {
		Payload struct {
			Binding struct {
				UUID string
			}
		} `graphql:"removeBinding(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
}
