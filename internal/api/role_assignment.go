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
	Principal RoleAssignmentPrincipal
	Target    RoleAssignmentTarget
}

// RoleAssignmentPrincipal represents the principal (user or SSO group) for a role assignment.
type RoleAssignmentPrincipal struct {
	Type string `json:"type"` // "user" or "sso-group"
	ID   int64  `json:"id"`
}

// RoleAssignmentTarget represents the target entity for a role assignment.
type RoleAssignmentTarget struct {
	Type string  `json:"type"`           // "system", "account-group", "policy-collection", "repository"
	UUID *string `json:"uuid,omitempty"` // Required for all target types except "system"
}

// RoleAssignmentInput is the input for creating a role assignment.
type RoleAssignmentInput struct {
	RoleName  string                  `json:"roleName"`
	Principal RoleAssignmentPrincipal `json:"principal"`
	Target    RoleAssignmentTarget    `json:"target"`
}

func (i RoleAssignmentInput) GetGraphQLType() string {
	return "UpdateRoleAssignmentInput"
}

type roleAssignmentAPI struct {
	c *graphql.Client
}

// Create creates a role assignment.
func (r roleAssignmentAPI) Create(ctx context.Context, input RoleAssignmentInput) (*RoleAssignment, error) {
	var mutation struct {
		Payload struct {
			RoleAssignments []RoleAssignment
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	variables := map[string]any{"input": input}
	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	// The mutation returns a list of role assignments, find the one we just created
	for _, assignment := range mutation.Payload.RoleAssignments {
		if assignment.Role.Name == input.RoleName &&
			assignment.Principal.Type == input.Principal.Type &&
			assignment.Principal.ID == input.Principal.ID {
			return &assignment, nil
		}
	}

	return nil, APIError{Kind: "API Error", Detail: "Role assignment not found in response"}
}

// Read returns data for a specific role assignment by ID.
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

// List returns role assignments filtered by target or principal.
func (r roleAssignmentAPI) List(ctx context.Context, target *RoleAssignmentTarget, principal *RoleAssignmentPrincipal) ([]RoleAssignment, error) {
	var query struct {
		RoleAssignments struct {
			Edges []struct {
				Node RoleAssignment
			}
		} `graphql:"roleAssignments(target: $target, principal: $principal)"`
	}

	variables := map[string]any{
		"target":    target,
		"principal": principal,
	}

	if err := r.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	assignments := make([]RoleAssignment, 0, len(query.RoleAssignments.Edges))
	for _, edge := range query.RoleAssignments.Edges {
		assignments = append(assignments, edge.Node)
	}

	return assignments, nil
}

// Delete removes a role assignment.
func (r roleAssignmentAPI) Delete(ctx context.Context, input RoleAssignmentInput) error {
	// To delete a role assignment, we call updateRoleAssignment with an empty role name
	deleteInput := RoleAssignmentInput{
		RoleName:  "", // Empty role name indicates removal
		Principal: input.Principal,
		Target:    input.Target,
	}

	var mutation struct {
		Payload struct {
			RoleAssignments []RoleAssignment
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	variables := map[string]any{"input": deleteInput}
	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}

	return nil
}
