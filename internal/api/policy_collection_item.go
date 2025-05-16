package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// PolicyCollectionItem is the data returned by reading a policy collection mapping.
type PolicyCollectionItem struct {
	ID     string
	Policy struct {
		UUID    string
		Version int
	}
	Collection struct {
		UUID string
	}
}

// PolicyCollectionItemCreateInput is the input for creating a policy collection item.
type PolicyCollectionItemCreateInput struct {
	CollectionUUID string `json:"collectionUUID"`
	PolicyUUID     string `json:"policyUUID"`
	PolicyVersion  int    `json:"policyVersion"`
}

func (i PolicyCollectionItemCreateInput) GetGraphQLType() string {
	return "PolicyCollectionMappingsInput"
}

type upsertPolicyCollectionMappingsInput struct {
	Mappings []PolicyCollectionItemCreateInput `json:"mappings"`
}

func (i upsertPolicyCollectionMappingsInput) GetGraphQLType() string {
	return "UpsertPolicyCollectionMappingsInput"
}

type removePolicyCollectionItemInput struct {
	IDs []graphql.ID `json:"ids"`
}

func (i removePolicyCollectionItemInput) GetGraphQLType() string {
	return "RemovePolicyCollectionMappingsInput"
}

type policyCollectionItemAPI struct {
	c *graphql.Client
}

// Read returns data for a policy collection item.
func (a policyCollectionItemAPI) Read(ctx context.Context, collectionUUID string, policyUUID string) (PolicyCollectionItem, error) {
	var query struct {
		PolicyCollection struct {
			PolicyMappings struct {
				Edges []struct {
					Node PolicyCollectionItem
				}
			}
		} `graphql:"policyCollection(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(collectionUUID),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return PolicyCollectionItem{}, APIError{"Client error", err.Error()}
	}

	for _, edge := range query.PolicyCollection.PolicyMappings.Edges {
		if edge.Node.Policy.UUID == policyUUID {
			return edge.Node, nil
		}
	}

	return PolicyCollectionItem{}, nil
}

// Create creates a policy collection item.
func (a policyCollectionItemAPI) Create(ctx context.Context, input PolicyCollectionItemCreateInput) (PolicyCollectionItem, error) {
	var mutation struct {
		UpsertPolicyCollectionMappings struct {
			Mappings []PolicyCollectionItem
		} `graphql:"upsertPolicyCollectionMappings(input: $input)"`
	}
	variables := map[string]any{
		"input": upsertPolicyCollectionMappingsInput{
			Mappings: []PolicyCollectionItemCreateInput{input},
		},
	}

	err := a.c.Mutate(ctx, &mutation, variables)
	if err != nil {
		return PolicyCollectionItem{}, APIError{"Client error", err.Error()}
	}

	if len(mutation.UpsertPolicyCollectionMappings.Mappings) == 0 {
		return PolicyCollectionItem{}, APIError{"Create error", "Policy collection item not found after creation"}
	}

	return mutation.UpsertPolicyCollectionMappings.Mappings[0], nil
}

// Delete removes a policy collection item.
func (a policyCollectionItemAPI) Delete(ctx context.Context, id string) error {
	var mutation struct {
		RemovePolicyCollectionMappings struct {
			Removed []struct {
				ID string
			}
		} `graphql:"removePolicyCollectionMappings(input: $input)"`
	}
	variables := map[string]any{
		"input": removePolicyCollectionItemInput{
			IDs: []graphql.ID{graphql.ID(id)},
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return APIError{"Client error", err.Error()}
	}
	return nil
}
