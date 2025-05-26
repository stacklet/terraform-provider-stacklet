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
	ExecutionConfig struct {
		Variables *string
	}
	System bool
}

// ExecutionConfig holds the execution configuration for a binding.
type BindingExecutionConfig struct {
	Variables *string `json:"variables"`
}

// BindingCreateInput is the input for creating a binding.
type BindingCreateInput struct {
	Name                 string                 `json:"name"`
	Description          *string                `json:"description,omitempty"`
	AutoDeploy           bool                   `json:"autoDeploy"`
	Schedule             *string                `json:"schedule,omitempty"`
	ExecutionConfig      BindingExecutionConfig `json:"executionConfig,omitempty"`
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
func (a bindingAPI) Read(ctx context.Context, uuid string, name string) (Binding, error) {
	var query struct {
		Binding Binding `graphql:"binding(uuid: $uuid, name: $name)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return query.Binding, APIError{"Client error", err.Error()}
	}
	if query.Binding.ID == "" {
		return query.Binding, NotFound{"Binding not found"}
	}

	return query.Binding, nil
}

// Create creates a binding.
func (a bindingAPI) Create(ctx context.Context, i BindingCreateInput) (Binding, error) {
	var mutation struct {
		AddBinding struct {
			Binding Binding
		} `graphql:"addBinding(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return mutation.AddBinding.Binding, APIError{"Client error", err.Error()}
	}

	return mutation.AddBinding.Binding, nil
}

// Update updates a binding.
func (a bindingAPI) Update(ctx context.Context, i BindingUpdateInput) (Binding, error) {
	var mutation struct {
		UpdateBinding struct {
			Binding Binding
		} `graphql:"updateBinding(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return mutation.UpdateBinding.Binding, APIError{"Client error", err.Error()}
	}

	return mutation.UpdateBinding.Binding, nil
}

// Delete removes a binding.
func (a bindingAPI) Delete(ctx context.Context, uuid string) error {
	var mutation struct {
		RemoveBinding struct {
			Binding struct {
				UUID string
			}
		} `graphql:"removeBinding(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return APIError{"Client error", err.Error()}
	}
	return nil
}
