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
}

// PolicyCollectionCreateInput is the input to create a policy collection.
type PolicyCollectionCreateInput struct {
	Name        string        `json:"name"`
	Provider    CloudProvider `json:"provider"`
	Description *string       `json:"description,omitempty"`
	AutoUpdate  *bool         `json:"autoUpdate,omitempty"`
}

func (i PolicyCollectionCreateInput) GetGraphQLType() string {
	return "AddPolicyCollectionInput"
}

type PolicyCollectionUpdateInput struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	Provider    CloudProvider `json:"provider"`
	Description *string       `json:"description"`
	AutoUpdate  *bool         `json:"autoUpdate"`
}

func (i PolicyCollectionUpdateInput) GetGraphQLType() string {
	return "UpdatePolicyCollectionInput"
}

type policyCollectionAPI struct {
	c *graphql.Client
}

// Read returns data for an account.
func (a policyCollectionAPI) Read(ctx context.Context, uuid string, name string) (PolicyCollection, error) {
	var query struct {
		PolicyCollection PolicyCollection `graphql:"policyCollection(uuid: $uuid, name: $name)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return query.PolicyCollection, APIError{"Client error", err.Error()}
	}

	return query.PolicyCollection, nil
}

// Create creates a policy collection.
func (a policyCollectionAPI) Create(ctx context.Context, i PolicyCollectionCreateInput) (PolicyCollection, error) {
	var mutation struct {
		AddPolicyCollection struct {
			Collection PolicyCollection
		} `graphql:"addPolicyCollection(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return mutation.AddPolicyCollection.Collection, APIError{"Client error", err.Error()}
	}
	return mutation.AddPolicyCollection.Collection, nil
}

// Update updates a policy collection.
func (a policyCollectionAPI) Update(ctx context.Context, i PolicyCollectionUpdateInput) (PolicyCollection, error) {
	var mutation struct {
		UpdatePolicyCollection struct {
			Collection PolicyCollection
		} `graphql:"updatePolicyCollection(input: $input)"`
	}
	input := map[string]any{"input": i}
	if err := a.c.Mutate(ctx, &mutation, input); err != nil {
		return mutation.UpdatePolicyCollection.Collection, APIError{"Client error", err.Error()}
	}

	return mutation.UpdatePolicyCollection.Collection, nil
}

// Delete removes a policy collection.
func (a policyCollectionAPI) Delete(ctx context.Context, uuid string) error {
	var mutation struct {
		RemovePolicyCollection struct {
			Collection struct {
				UUID string
			}
		} `graphql:"removePolicyCollection(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return APIError{"Client error", err.Error()}
	}
	return nil
}
