// Copyright (c) 2025 - Stacklet, Inc.

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
func (a policyCollectionMappingAPI) Read(ctx context.Context, collectionUUID string, policyUUID string) (*PolicyCollectionMapping, error) {
	cursor := ""
	var query struct {
		PolicyCollection struct {
			PolicyMappings struct {
				Edges []struct {
					Node PolicyCollectionMapping
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
				Problems []Problem
			} `graphql:"policyMappings(first: 100, after: $cursor)"`
		} `graphql:"policyCollection(uuid: $uuid)"`
	}
	for {
		variables := map[string]any{
			"uuid":   graphql.String(collectionUUID),
			"cursor": graphql.String(cursor),
		}
		if err := a.c.Query(ctx, &query, variables); err != nil {
			return nil, NewAPIError(err)
		}
		if err := FromProblems(ctx, query.PolicyCollection.PolicyMappings.Problems); err != nil {
			return nil, err
		}

		for _, edge := range query.PolicyCollection.PolicyMappings.Edges {
			if edge.Node.Policy.UUID == policyUUID {
				return &edge.Node, nil
			}
		}

		if !query.PolicyCollection.PolicyMappings.PageInfo.HasNextPage {
			return nil, NotFound{"Policy collection mapping not found"}
		}
		cursor = query.PolicyCollection.PolicyMappings.PageInfo.EndCursor
	}
}

// Upsert creates or updates a policy collection mapping.
func (a policyCollectionMappingAPI) Upsert(ctx context.Context, input PolicyCollectionMappingInput) (*PolicyCollectionMapping, error) {
	var mutation struct {
		Payload struct {
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
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Mappings[0], nil
}

// Delete removes a policy collection mapping.
func (a policyCollectionMappingAPI) Delete(ctx context.Context, id string) error {
	var mutation struct {
		Payload struct {
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
		return NewAPIError(err)
	}
	return nil
}
