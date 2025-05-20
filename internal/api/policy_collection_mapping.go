package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// PolicyCollectionMapping is the data returned by reading a policy collection mapping.
type PolicyCollectionMapping struct {
	ID     string
	Policy struct {
		UUID    string
		Version int
	}
	Collection struct {
		UUID string
	}
}

// PolicyCollectionMappingInput is the input for creating or updating a policy collection mapping.
type PolicyCollectionMappingInput struct {
	CollectionUUID string `json:"collectionUUID"`
	PolicyUUID     string `json:"policyUUID"`
	PolicyVersion  int    `json:"policyVersion"`
}

func (i PolicyCollectionMappingInput) GetGraphQLType() string {
	return "PolicyCollectionMappingsInput"
}

type upsertPolicyCollectionMappingsInput struct {
	Mappings []PolicyCollectionMappingInput `json:"mappings"`
}

func (i upsertPolicyCollectionMappingsInput) GetGraphQLType() string {
	return "UpsertPolicyCollectionMappingsInput"
}

type removePolicyCollectionMappingInput struct {
	IDs []graphql.ID `json:"ids"`
}

func (i removePolicyCollectionMappingInput) GetGraphQLType() string {
	return "RemovePolicyCollectionMappingsInput"
}

type policyCollectionMappingAPI struct {
	c *graphql.Client
}

// Read returns data for a policy collection mapping.
func (a policyCollectionMappingAPI) Read(ctx context.Context, collectionUUID string, policyUUID string) (PolicyCollectionMapping, error) {
	var query struct {
		PolicyCollection struct {
			PolicyMappings struct {
				Edges []struct {
					Node PolicyCollectionMapping
				}
			}
		} `graphql:"policyCollection(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(collectionUUID),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return PolicyCollectionMapping{}, APIError{"Client error", err.Error()}
	}

	for _, edge := range query.PolicyCollection.PolicyMappings.Edges {
		if edge.Node.Policy.UUID == policyUUID {
			return edge.Node, nil
		}
	}

	return PolicyCollectionMapping{}, nil
}

// Create creates a policy collection mapping.
func (a policyCollectionMappingAPI) Create(ctx context.Context, input PolicyCollectionMappingInput) (PolicyCollectionMapping, error) {
	return a.upsertMapping(ctx, input)
}

// Update updates a policy collection mapping.
func (a policyCollectionMappingAPI) Update(ctx context.Context, input PolicyCollectionMappingInput) (PolicyCollectionMapping, error) {
	return a.upsertMapping(ctx, input)
}

// Delete removes a policy collection mapping.
func (a policyCollectionMappingAPI) Delete(ctx context.Context, id string) error {
	var mutation struct {
		RemovePolicyCollectionMappings struct {
			Removed []struct {
				ID string
			}
		} `graphql:"removePolicyCollectionMappings(input: $input)"`
	}
	variables := map[string]any{
		"input": removePolicyCollectionMappingInput{
			IDs: []graphql.ID{graphql.ID(id)},
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return APIError{"Client error", err.Error()}
	}
	return nil
}

func (a policyCollectionMappingAPI) upsertMapping(ctx context.Context, input PolicyCollectionMappingInput) (PolicyCollectionMapping, error) {
	var mutation struct {
		UpsertPolicyCollectionMappings struct {
			Mappings []PolicyCollectionMapping
		} `graphql:"upsertPolicyCollectionMappings(input: $input)"`
	}
	variables := map[string]any{
		"input": upsertPolicyCollectionMappingsInput{
			Mappings: []PolicyCollectionMappingInput{input},
		},
	}

	err := a.c.Mutate(ctx, &mutation, variables)
	if err != nil {
		return PolicyCollectionMapping{}, APIError{"Client error", err.Error()}
	}

	if len(mutation.UpsertPolicyCollectionMappings.Mappings) == 0 {
		return PolicyCollectionMapping{}, APIError{"Create error", "Policy collection mapping not found after creation"}
	}

	return mutation.UpsertPolicyCollectionMappings.Mappings[0], nil
}
