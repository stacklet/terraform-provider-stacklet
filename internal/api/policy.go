package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Policy is the data returned by reading policy data
type Policy struct {
	ID          string
	UUID        string
	Name        string
	Description *string
	Provider    string
	Version     float64
}

type policyAPI struct {
	c *graphql.Client
}

// Read returns data for a policy
func (a policyAPI) Read(ctx context.Context, uuid string, name string) (Policy, error) {
	var query struct {
		Policy Policy `graphql:"policy(uuid: $uuid, name: $name)"`
	}
	variables := map[string]any{
		"uuid": graphql.String(uuid),
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return query.Policy, APIError{"Client error", err.Error()}
	}

	return query.Policy, nil
}
