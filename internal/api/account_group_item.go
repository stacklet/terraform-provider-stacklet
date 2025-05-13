package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// AccountGroupItem is the data returned by reading an account group item data.
type AccountGroupItem struct {
	ID         string
	GroupUUID  string
	AccountKey string
	Provider   CloudProvider
}

type upsertAccountGroupMappingsInput struct {
	Mappings []accountGroupMappingInput `json:"mappings"`
}

func (i upsertAccountGroupMappingsInput) GetGraphQLType() string {
	return "UpsertAccountGroupMappingsInput"
}

type accountGroupMappingInput struct {
	AccountKey string `json:"accountKey"`
	GroupUUID  string `json:"groupUUID"`
}

func (i accountGroupMappingInput) GetGraphQLType() string {
	return "AccountGroupMappingInput"
}

type removeAccountGroupMappingsInput struct {
	IDs []graphql.ID `json:"ids"`
}

func (i removeAccountGroupMappingsInput) GetGraphQLType() string {
	return "RemoveAccountGroupMappingsInput"
}

type accountGroupItemAPI struct {
	c *graphql.Client
}

// Read returns data for an account group item.
func (a accountGroupItemAPI) Read(ctx context.Context, cloudProvider string, accountKey string, groupUUID string) (AccountGroupItem, error) {
	provider, err := NewCloudProvider(cloudProvider)
	if err != nil {
		return AccountGroupItem{}, APIError{"Invalid provider", err.Error()}
	}

	var query struct {
		AccountGroup struct {
			AccountMappings struct {
				Edges []struct {
					Node struct {
						ID      string
						Account struct {
							Key      string
							Provider string
						}
					}
				}
			}
		} `graphql:"accountGroup(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(groupUUID),
	}
	if err = a.c.Query(ctx, &query, variables); err != nil {
		return AccountGroupItem{}, APIError{"Client error", err.Error()}
	}

	for _, edge := range query.AccountGroup.AccountMappings.Edges {
		if edge.Node.Account.Key == accountKey && edge.Node.Account.Provider == cloudProvider {
			return AccountGroupItem{
				ID:         edge.Node.ID,
				GroupUUID:  groupUUID,
				AccountKey: accountKey,
				Provider:   provider,
			}, nil
		}
	}

	return AccountGroupItem{}, nil
}

// Create creates an account group item.
func (a accountGroupItemAPI) Create(ctx context.Context, cloudProvider string, accountKey string, groupUUID string) (AccountGroupItem, error) {
	provider, err := NewCloudProvider(cloudProvider)
	if err != nil {
		return AccountGroupItem{}, APIError{"Invalid provider", err.Error()}
	}

	var mutation struct {
		UpsertAccountGroupMappings struct {
			Mappings []struct {
				ID string
			}
		} `graphql:"upsertAccountGroupMappings(input: $input)"`
	}
	variables := map[string]any{
		"input": upsertAccountGroupMappingsInput{
			Mappings: []accountGroupMappingInput{
				{
					AccountKey: accountKey,
					GroupUUID:  groupUUID,
				},
			},
		},
	}

	err = a.c.Mutate(ctx, &mutation, variables)
	if err != nil {
		return AccountGroupItem{}, APIError{"Client error", err.Error()}
	}

	return AccountGroupItem{
		ID:         mutation.UpsertAccountGroupMappings.Mappings[0].ID,
		AccountKey: accountKey,
		GroupUUID:  groupUUID,
		Provider:   provider,
	}, nil
}

// Delete removes an account group item.
func (a accountGroupItemAPI) Delete(ctx context.Context, id string) error {
	var mutation struct {
		RemoveAccountGroupMappings struct {
			Removed []struct {
				ID string
			}
		} `graphql:"removeAccountGroupMappings(input: $input)"`
	}
	variables := map[string]any{
		"input": removeAccountGroupMappingsInput{
			IDs: []graphql.ID{graphql.ID(id)},
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return APIError{"Client error", err.Error()}
	}
	return nil
}
