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
	RoleAssignmentTarget string          `graphql:"roleAssignmentTarget"`
	RoleScope            *roleTargetType `graphql:"... on RoleScope"`
	AccountGroup         *roleTargetType `graphql:"... on AccountGroup"`
	PolicyCollection     *roleTargetType `graphql:"... on PolicyCollection"`
	Repository           *roleTargetType `graphql:"... on Repository"`
	RepositoryConfig     *roleTargetType `graphql:"... on RepositoryConfig"`
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

// Create assigns a role to a principal on a target.
// roleName is the name of the role to assign.
// principal and target are opaque string identifiers.
func (r roleAssignmentAPI) Create(ctx context.Context, roleName string, principal string, target string) (*RoleAssignment, error) {
	var mutation struct {
		Payload struct {
			Assignment RoleAssignment
		} `graphql:"assignRole(roleName: $roleName, principal: $principal, target: $target)"`
	}
	variables := map[string]any{
		"roleName":  graphql.String(roleName),
		"principal": graphql.String(principal),
		"target":    graphql.String(target),
	}
	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Assignment, nil
}

// Read returns a single role assignment by ID.
func (r roleAssignmentAPI) Read(ctx context.Context, id string) (*RoleAssignment, error) {
	var query struct {
		RoleAssignment RoleAssignment `graphql:"roleAssignment(id: $id)"`
	}

	variables := map[string]any{
		"id": graphql.String(id),
	}
	if err := r.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if query.RoleAssignment.ID == "" {
		return nil, NotFound{"Role assignment not found"}
	}

	return &query.RoleAssignment, nil
}

// Delete removes a role assignment.
// roleName is the name of the role to unassign.
// principal and target are opaque string identifiers.
func (r roleAssignmentAPI) Delete(ctx context.Context, roleName string, principal string, target string) error {
	var mutation struct {
		Payload struct {
			Success bool
		} `graphql:"unassignRole(roleName: $roleName, principal: $principal, target: $target)"`
	}
	variables := map[string]any{
		"roleName":  graphql.String(roleName),
		"principal": graphql.String(principal),
		"target":    graphql.String(target),
	}
	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
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
