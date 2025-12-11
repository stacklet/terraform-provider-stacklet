// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// RoleAssignment is the data returned by reading role assignment data.
type RoleAssignment struct {
	ID        string
	Role      Role
	Principal rolePrincipal `graphql:"principal"`
	Target    roleTarget    `graphql:"target"`
}

// rolePrincipalPrincipal contains the opaque roleAssignmentPrincipal string.
type rolePrincipalPrincipal struct {
	RoleAssignmentPrincipal string `graphql:"roleAssignmentPrincipal"`
}

// rolePrincipal represents the GraphQL union type for RolePrincipal.
type rolePrincipal struct {
	User     *rolePrincipalPrincipal `graphql:"... on User"`
	SSOGroup *rolePrincipalPrincipal `graphql:"... on SSOGroup"`
}

// GetPrincipal extracts the opaque principal identifier string.
func (r *RoleAssignment) GetPrincipal() string {
	if r.Principal.User != nil {
		return r.Principal.User.RoleAssignmentPrincipal
	}
	if r.Principal.SSOGroup != nil {
		return r.Principal.SSOGroup.RoleAssignmentPrincipal
	}
	return ""
}

// roleTarget represents the GraphQL union type for target entities.
type roleTarget struct {
	RoleAssignmentTarget string                  `graphql:"roleAssignmentTarget"`
	RoleScope            *roleTargetType         `graphql:"... on RoleScope"`
	AccountGroup         *roleTargetType         `graphql:"... on AccountGroup"`
	PolicyCollection     *roleTargetType         `graphql:"... on PolicyCollection"`
	Repository           *roleTargetType         `graphql:"... on Repository"`
	RepositoryConfig     *roleTargetType         `graphql:"... on RepositoryConfig"`
}

// roleTargetType is used for the union type matching.
type roleTargetType struct {
	RoleAssignmentTarget string `graphql:"roleAssignmentTarget"`
}

// GetTarget extracts the opaque target identifier string.
func (r *RoleAssignment) GetTarget() string {
	return r.Target.RoleAssignmentTarget
}

type roleAssignmentAPI struct {
	c *graphql.Client
}

// List returns role assignments, optionally filtered by target or principal.
// target and principal are opaque string identifiers. Pass nil to skip filtering.
func (r roleAssignmentAPI) List(ctx context.Context, target *string, principal *string) ([]RoleAssignment, error) {
	cursor := ""
	var query struct {
		RoleAssignments struct {
			Edges []struct {
				Node RoleAssignment
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"roleAssignments(first: 100, after: $cursor)"`
	}

	assignments := make([]RoleAssignment, 0)

	// Paginate through all results
	for {
		variables := map[string]any{
			"cursor": graphql.String(cursor),
		}

		if err := r.c.Query(ctx, &query, variables); err != nil {
			return nil, NewAPIError(err)
		}

		for _, edge := range query.RoleAssignments.Edges {
			assignment := edge.Node

			// Filter by target if specified (client-side filtering)
			if target != nil && assignment.GetTarget() != *target {
				continue
			}

			// Filter by principal if specified (client-side filtering)
			if principal != nil && assignment.GetPrincipal() != *principal {
				continue
			}

			assignments = append(assignments, assignment)
		}

		// Check if there are more pages
		if !query.RoleAssignments.PageInfo.HasNextPage {
			break
		}
		cursor = query.RoleAssignments.PageInfo.EndCursor
	}

	return assignments, nil
}

// ListByTargetString returns role assignments filtered by an opaque target string.
// The target string should be in the format "type:id" (e.g., "account-group:uuid", "system:all").
func (r roleAssignmentAPI) ListByTargetString(ctx context.Context, targetStr string) ([]RoleAssignment, error) {
	// Get all role assignments (no server-side filtering for now)
	// We'll filter client-side by comparing the target strings
	assignments, err := r.List(ctx, &targetStr, nil)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}
