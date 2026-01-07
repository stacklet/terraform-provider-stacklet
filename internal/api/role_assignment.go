// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"
	"fmt"

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

// RoleAssignmentInput represents the input for granting or revoking a role assignment.
type RoleAssignmentInput struct {
	RoleName  string `json:"roleName"`
	Principal string `json:"principal"`
	Target    string `json:"target"`
}

// UpdateRoleAssignmentInput represents the input for the updateRoleAssignment mutation.
type UpdateRoleAssignmentInput struct {
	Grant  []RoleAssignmentInput `json:"grant,omitempty"`
	Revoke []RoleAssignmentInput `json:"revoke,omitempty"`
}

// GrantRoleAssignmentPayload represents the result of granting a role assignment.
type GrantRoleAssignmentPayload struct {
	ErrorMessage   *string
	RoleAssignment *RoleAssignment
}

func (p GrantRoleAssignmentPayload) Error() string {
	if p.ErrorMessage == nil {
		return ""
	}

	return *p.ErrorMessage
}

// RevokeRoleAssignmentPayload represents the result of revoking a role assignment.
type RevokeRoleAssignmentPayload struct {
	ErrorMessage *string
	Removed      struct {
		ID string
	}
}

func (p RevokeRoleAssignmentPayload) Error() string {
	if p.ErrorMessage == nil {
		return ""
	}

	return *p.ErrorMessage
}

// Create assigns a role to a principal on a target.
// roleName is the name of the role to assign.
// principal and target are opaque string identifiers.
func (r roleAssignmentAPI) Create(ctx context.Context, roleName string, principal string, target string) (*RoleAssignment, error) {
	// Use updateRoleAssignment mutation with grant list
	var mutation struct {
		UpdateRoleAssignment struct {
			Grant  []GrantRoleAssignmentPayload
			Revoke []RevokeRoleAssignmentPayload
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	input := UpdateRoleAssignmentInput{
		Grant: []RoleAssignmentInput{
			{
				RoleName:  roleName,
				Principal: principal,
				Target:    target,
			},
		},
	}

	variables := map[string]any{
		"input": input,
	}

	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	// Check if we have a grant result
	if len(mutation.UpdateRoleAssignment.Grant) > 0 {
		grantPayload := mutation.UpdateRoleAssignment.Grant[0]

		// Check for errors
		if grantPayload.ErrorMessage != nil && *grantPayload.ErrorMessage != "" {
			return nil, NewAPIError(fmt.Errorf("failed to grant role assignment: %s", *grantPayload.ErrorMessage))
		}

		// The ID returned from the mutation may not match the actual persisted assignment
		// Query to get the actual role assignment by the unique (role, principal, target) combination
		if grantPayload.RoleAssignment != nil {
			// List assignments filtered by target and principal, then find the matching role
			assignments, err := r.List(ctx, &target, &principal)
			if err != nil {
				return nil, err
			}

			// Find the assignment with the matching role name
			for _, assignment := range assignments {
				if assignment.Role.Name == roleName {
					return &assignment, nil
				}
			}
		}
	}

	return nil, NotFound{"Role assignment not found after creation"}
}

// Read returns a single role assignment by the unique combination of roleName, principal, and target.
// The roleAssignments API doesn't support filtering by ID, so we use the composite key to identify the assignment.
func (r roleAssignmentAPI) Read(ctx context.Context, roleName string, principal string, target string) (*RoleAssignment, error) {
	// Fetch role assignments filtered by principal and target
	assignments, err := r.List(ctx, &target, &principal)
	if err != nil {
		return nil, err
	}

	// Find the assignment with the matching role name
	for _, assignment := range assignments {
		if assignment.Role.Name == roleName {
			return &assignment, nil
		}
	}

	return nil, NotFound{"Role assignment not found"}
}

// Delete removes a role assignment.
// roleName is the name of the role to unassign.
// principal and target are opaque string identifiers.
func (r roleAssignmentAPI) Delete(ctx context.Context, roleName string, principal string, target string) error {
	// Use updateRoleAssignment mutation with revoke list
	var mutation struct {
		UpdateRoleAssignment struct {
			Grant  []GrantRoleAssignmentPayload
			Revoke []RevokeRoleAssignmentPayload
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	input := UpdateRoleAssignmentInput{
		Revoke: []RoleAssignmentInput{
			{
				RoleName:  roleName,
				Principal: principal,
				Target:    target,
			},
		},
	}

	variables := map[string]any{
		"input": input,
	}

	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}

	// Check if we have a revoke result with an error
	if len(mutation.UpdateRoleAssignment.Revoke) > 0 {
		revokePayload := mutation.UpdateRoleAssignment.Revoke[0]
		if revokePayload.ErrorMessage != nil && *revokePayload.ErrorMessage != "" {
			return NewAPIError(fmt.Errorf("failed to revoke role assignment: %s", *revokePayload.ErrorMessage))
		}
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
