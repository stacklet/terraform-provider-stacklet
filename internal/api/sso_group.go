// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// SSOGroup is the data returned by reading SSO group data.
type SSOGroup struct {
	ID                      string
	DisplayName             *string
	Name                    string
	RoleAssignmentPrincipal string
}

type ssoGroupAPI struct {
	c *graphql.Client
}

// Read returns data for an SSO group by name.
// Note: The Stacklet API SSO group filter requires the field name "name" with no operator.
func (s ssoGroupAPI) Read(ctx context.Context, name string) (*SSOGroup, error) {
	var query struct {
		SSOGroups struct {
			Edges []struct {
				Node SSOGroup
			}
		} `graphql:"ssoGroups(filterElement: $filterElement)"`
	}
	// Use "name" as field name and omit operator
	variables := map[string]any{
		"filterElement": newSimpleFilter("name", name),
	}
	if err := s.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(query.SSOGroups.Edges) == 0 {
		return nil, NotFound{"SSO group not found"}
	}

	return &query.SSOGroups.Edges[0].Node, nil
}
