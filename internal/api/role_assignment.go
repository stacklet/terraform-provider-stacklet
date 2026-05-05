// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

// RoleAssignment is the data returned by reading role assignment data.
type RoleAssignment struct {
	ID        graphql.ID
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
	c *client
}

type roleAssignmentItemInput struct {
	RoleName  string `json:"roleName"`
	Principal string `json:"principal"`
	Target    string `json:"target"`
}

type roleAssignmentInput struct {
	Grant  []roleAssignmentItemInput `json:"grant,omitempty"`
	Revoke []roleAssignmentItemInput `json:"revoke,omitempty"`
}

func (i roleAssignmentInput) GetGraphQLType() string {
	return "UpdateRoleAssignmentInput"
}

// grantRoleAssignmentPayload represents the result of granting a role assignment.
type grantRoleAssignmentPayload struct {
	ErrorMessage   *string
	RoleAssignment *RoleAssignment
}

func (p grantRoleAssignmentPayload) Error() string {
	if p.ErrorMessage == nil {
		return ""
	}

	return *p.ErrorMessage
}

// revokeRoleAssignmentPayload represents the result of revoking a role assignment.
type revokeRoleAssignmentPayload struct {
	ErrorMessage *string
	Removed      struct {
		ID graphql.ID
	}
}

func (p revokeRoleAssignmentPayload) Error() string {
	if p.ErrorMessage == nil {
		return ""
	}

	return *p.ErrorMessage
}

// Create assigns a role to a principal on a target.
// roleName is the name of the role to assign.
// principal and target are opaque string identifiers.
func (r roleAssignmentAPI) Create(ctx context.Context, roleName string, principal string, target string) (*RoleAssignment, error) {
	var mutation struct {
		UpdateRoleAssignment struct {
			Grant []grantRoleAssignmentPayload
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	input := roleAssignmentInput{
		Grant: []roleAssignmentItemInput{
			{
				RoleName:  roleName,
				Principal: principal,
				Target:    target,
			},
		},
	}

	if err := r.c.Mutate(ctx, &mutation, map[string]any{"input": input}); err != nil {
		return nil, err
	}

	if len(mutation.UpdateRoleAssignment.Grant) == 0 {
		return nil, NotFound{"Role assignment not found after creation"}
	}

	grantPayload := mutation.UpdateRoleAssignment.Grant[0]
	if grantPayload.Error() != "" {
		return nil, newAPIError(fmt.Errorf("failed to grant role assignment: %w", grantPayload))
	}

	if grantPayload.RoleAssignment == nil {
		return nil, NotFound{"Role assignment not found after creation"}
	}

	return r.Read(ctx, roleName, principal, target)
}

// Read returns a single role assignment by the unique combination of roleName, principal, and target.
// The roleAssignments API doesn't support filtering by ID, so we use the composite key to identify the assignment.
func (r roleAssignmentAPI) Read(ctx context.Context, roleName string, principal string, target string) (*RoleAssignment, error) {
	filters := []filterElementInput{
		newExactMatchFilter("role-name", roleName),
		newExactMatchFilter("target", target),
	}
	assignments, err := r.list(ctx, newCompositeFilter(filters, filterBooleanAND))
	if err != nil {
		return nil, err
	}

	for _, assignment := range assignments {
		if assignment.GetPrincipal() == principal {
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
			Revoke []revokeRoleAssignmentPayload
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	input := roleAssignmentInput{
		Revoke: []roleAssignmentItemInput{
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
		return err
	}

	// Check if we have a revoke result with an error
	if len(mutation.UpdateRoleAssignment.Revoke) > 0 {
		revokePayload := mutation.UpdateRoleAssignment.Revoke[0]
		if revokePayload.Error() != "" {
			return newAPIError(fmt.Errorf("failed to revoke role assignment: %w", revokePayload))
		}
	}

	return nil
}

// List returns role assignments for a target.
func (r roleAssignmentAPI) List(ctx context.Context, target string) ([]RoleAssignment, error) {
	return r.list(ctx, newExactMatchFilter("target", target))
}

func (r roleAssignmentAPI) list(ctx context.Context, filter filterElementInput) ([]RoleAssignment, error) {
	cursor := ""
	assignments := make([]RoleAssignment, 0)
	for {
		var query struct {
			RoleAssignments struct {
				Edges []struct {
					Node RoleAssignment
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"roleAssignments(first: 100, after: $cursor, filterElement: $filterElement)"`
		}
		variables := map[string]any{
			"cursor":        graphql.String(cursor),
			"filterElement": filter,
		}

		if err := r.c.Query(ctx, &query, variables); err != nil {
			return nil, err
		}

		for _, edge := range query.RoleAssignments.Edges {
			assignments = append(assignments, edge.Node)
		}
		if !query.RoleAssignments.PageInfo.HasNextPage {
			break
		}
		cursor = query.RoleAssignments.PageInfo.EndCursor
	}

	return assignments, nil
}
