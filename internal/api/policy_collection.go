// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Account is the data returned by reading policy collection data.
type PolicyCollection struct {
	ID               string
	UUID             string
	Name             string
	Description      *string
	Provider         CloudProvider
	AutoUpdate       bool
	System           bool
	IsDynamic        bool
	RepositoryConfig struct {
		UUID *string
	}
	RepositoryView *RepositoryView
}

type RepositoryViewConfig struct {
	Namespace         string
	BranchName        string
	PolicyDirectories []string
	PolicyFileSuffix  []string
}

// PolicyCollectionCreateInput is the input to create a policy collection.
type PolicyCollectionCreateInput struct {
	Name           string               `json:"name"`
	Provider       CloudProvider        `json:"provider"`
	Description    *string              `json:"description,omitempty"`
	AutoUpdate     *bool                `json:"autoUpdate,omitempty"`
	RepositoryUUID *string              `json:"repositoryUUID,omitempty"`
	RepositoryView *RepositoryViewInput `json:"repositoryView,omitempty"`
}

func (i PolicyCollectionCreateInput) GetGraphQLType() string {
	return "AddPolicyCollectionInput"
}

type PolicyCollectionUpdateInput struct {
	UUID           string               `json:"uuid"`
	Name           string               `json:"name"`
	Provider       CloudProvider        `json:"provider"`
	Description    *string              `json:"description"`
	AutoUpdate     *bool                `json:"autoUpdate"`
	RepositoryView *RepositoryViewInput `json:"repositoryView,omitempty"`
}

func (i PolicyCollectionUpdateInput) GetGraphQLType() string {
	return "UpdatePolicyCollectionInput"
}

// RepositoryViewInput is the input for a repository view configuration in a policy collection.
type RepositoryViewInput struct {
	BranchName        *string  `json:"branchName"`
	PolicyDirectories []string `json:"policyDirectories"`
	PolicyFileSuffix  []string `json:"policyFileSuffix"`
}

type policyCollectionAPI struct {
	c *graphql.Client
}

// Read returns data for an account.
func (a policyCollectionAPI) Read(ctx context.Context, uuid string, name string) (*PolicyCollection, error) {
	var query struct {
		PolicyCollection PolicyCollection `graphql:"policyCollection(uuid: $uuid, name: $name)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if query.PolicyCollection.ID == "" {
		return nil, NotFound{"Policy collection not found"}
	}

	return &query.PolicyCollection, nil
}

// Create creates a policy collection.
func (a policyCollectionAPI) Create(ctx context.Context, i PolicyCollectionCreateInput) (*PolicyCollection, error) {
	var mutation struct {
		Payload struct {
			Collection PolicyCollection
		} `graphql:"addPolicyCollection(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}
	return &mutation.Payload.Collection, nil
}

// Update updates a policy collection.
func (a policyCollectionAPI) Update(ctx context.Context, i PolicyCollectionUpdateInput) (*PolicyCollection, error) {
	var mutation struct {
		Payload struct {
			Collection PolicyCollection
		} `graphql:"updatePolicyCollection(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Collection, nil
}

// Delete removes a policy collection.
func (a policyCollectionAPI) Delete(ctx context.Context, uuid string) error {
	var mutation struct {
		Payload struct {
			Collection struct {
				UUID string
			}
		} `graphql:"removePolicyCollection(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
}
