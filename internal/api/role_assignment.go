// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hasura/go-graphql-client"
)

// RoleAssignment is the data returned by reading role assignment data.
type RoleAssignment struct {
	ID        string
	Role      Role
	Principal rolePrincipal `graphql:"principal"`
	Target    roleTarget    `graphql:"target"`
}

// rolePrincipal represents the GraphQL union type for RolePrincipal.
type rolePrincipal struct {
	User     *principalUser     `graphql:"... on User"`
	SSOGroup *principalSSOGroup `graphql:"... on SSOGroup"`
}

// principalUser represents a User principal.
type principalUser struct {
	ID string `graphql:"id"`
}

// principalSSOGroup represents an SSOGroup principal.
type principalSSOGroup struct {
	ID string `graphql:"id"`
}

// decodeGraphQLNodeID decodes a GraphQL Node ID and extracts the numeric ID.
// GraphQL Node IDs are base64-encoded JSON arrays like ["user", "1"].
func decodeGraphQLNodeID(nodeID string) (int64, error) {
	decoded, err := base64.StdEncoding.DecodeString(nodeID)
	if err != nil {
		return 0, fmt.Errorf("failed to decode base64: %w", err)
	}

	var parts []interface{}
	if err := json.Unmarshal(decoded, &parts); err != nil {
		return 0, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid node ID format: expected at least 2 parts, got %d", len(parts))
	}

	// The second element is the numeric ID (could be string or number)
	switch v := parts[1].(type) {
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse ID as int64: %w", err)
		}
		return id, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unexpected ID type: %T", v)
	}
}

// roleTarget represents the GraphQL union type for target entities.
type roleTarget struct {
	Typename         string                  `graphql:"__typename"`
	RoleScope        *targetRoleScope        `graphql:"... on RoleScope"`
	AccountGroup     *targetAccountGroup     `graphql:"... on AccountGroup"`
	PolicyCollection *targetPolicyCollection `graphql:"... on PolicyCollection"`
	Repository       *targetRepository       `graphql:"... on Repository"`
	RepositoryConfig *targetRepositoryConfig `graphql:"... on RepositoryConfig"`
}

// targetRoleScope represents a system-level (RoleScope) target.
type targetRoleScope struct {
	Typename string `graphql:"__typename"`
}

// targetAccountGroup represents an AccountGroup target.
type targetAccountGroup struct {
	UUID string
}

// targetPolicyCollection represents a PolicyCollection target.
type targetPolicyCollection struct {
	UUID string
}

// targetRepository represents a Repository target.
type targetRepository struct {
	UUID string
}

// targetRepositoryConfig represents a RepositoryConfig target.
type targetRepositoryConfig struct {
	UUID string
}

// RoleAssignmentPrincipal represents the principal (user or SSO group) for a role assignment.
type RoleAssignmentPrincipal struct {
	Type string `json:"type"` // "user" or "sso-group"
	ID   int64  `json:"id"`
}

// String serializes the principal to the format expected by the GraphQL API.
func (p RoleAssignmentPrincipal) String() string {
	return fmt.Sprintf("%s:%d", p.Type, p.ID)
}

// GetPrincipal extracts the principal information from the GraphQL union type.
func (r *RoleAssignment) GetPrincipal() RoleAssignmentPrincipal {
	if r.Principal.User != nil {
		id, err := decodeGraphQLNodeID(r.Principal.User.ID)
		if err != nil {
			// If decoding fails, return empty principal
			return RoleAssignmentPrincipal{}
		}
		return RoleAssignmentPrincipal{
			Type: "user",
			ID:   id,
		}
	}
	if r.Principal.SSOGroup != nil {
		id, err := decodeGraphQLNodeID(r.Principal.SSOGroup.ID)
		if err != nil {
			// If decoding fails, return empty principal
			return RoleAssignmentPrincipal{}
		}
		return RoleAssignmentPrincipal{
			Type: "sso-group",
			ID:   id,
		}
	}
	return RoleAssignmentPrincipal{}
}

// RoleAssignmentTarget represents the target entity for a role assignment.
type RoleAssignmentTarget struct {
	Type string  `json:"type"`           // "system", "account-group", "policy-collection", "repository"
	UUID *string `json:"uuid,omitempty"` // Required for all target types except "system"
}

// String serializes the target to the format expected by the GraphQL API.
func (t RoleAssignmentTarget) String() string {
	if t.Type == "system" {
		return "system:all"
	}
	if t.UUID == nil {
		// Non-system targets require a UUID
		return ""
	}
	return fmt.Sprintf("%s:%s", t.Type, *t.UUID)
}

// GetTarget extracts the target information from the GraphQL union type.
func (r *RoleAssignment) GetTarget() RoleAssignmentTarget {
	if r.Target.RoleScope != nil {
		return RoleAssignmentTarget{
			Type: "system",
			UUID: nil,
		}
	}
	if r.Target.AccountGroup != nil {
		return RoleAssignmentTarget{
			Type: "account-group",
			UUID: &r.Target.AccountGroup.UUID,
		}
	}
	if r.Target.PolicyCollection != nil {
		return RoleAssignmentTarget{
			Type: "policy-collection",
			UUID: &r.Target.PolicyCollection.UUID,
		}
	}
	if r.Target.Repository != nil {
		return RoleAssignmentTarget{
			Type: "repository",
			UUID: &r.Target.Repository.UUID,
		}
	}
	if r.Target.RepositoryConfig != nil {
		return RoleAssignmentTarget{
			Type: "repository-config",
			UUID: &r.Target.RepositoryConfig.UUID,
		}
	}
	// Fallback to system if nothing matches
	return RoleAssignmentTarget{
		Type: "system",
		UUID: nil,
	}
}

// RoleAssignmentInput is the input for creating a role assignment.
type RoleAssignmentInput struct {
	Grants  []RoleAssignmentGrant  `json:"grant,omitempty"`
	Revokes []RoleAssignmentRevoke `json:"revoke,omitempty"`
}

// RoleAssignmentGrant represents a role to grant.
type RoleAssignmentGrant struct {
	Role      string `json:"roleName"`
	Principal string `json:"principal"`
	Target    string `json:"target"`
}

// RoleAssignmentRevoke represents a role to revoke.
type RoleAssignmentRevoke struct {
	Role      string `json:"roleName"`
	Principal string `json:"principal"`
	Target    string `json:"target"`
}

func (i RoleAssignmentInput) GetGraphQLType() string {
	return "UpdateRoleAssignmentInput"
}

type roleAssignmentAPI struct {
	c *graphql.Client
}

// Create creates a role assignment.
func (r roleAssignmentAPI) Create(ctx context.Context, roleName string, principal RoleAssignmentPrincipal, target RoleAssignmentTarget) (*RoleAssignment, error) {
	// Validate target
	targetStr := target.String()
	if targetStr == "" {
		return nil, APIError{
			Kind:   "Invalid Input",
			Detail: fmt.Sprintf("Target type '%s' requires a UUID, but none was provided", target.Type),
		}
	}

	input := RoleAssignmentInput{
		Grants: []RoleAssignmentGrant{
			{
				Role:      roleName,
				Principal: principal.String(),
				Target:    targetStr,
			},
		},
	}

	var mutation struct {
		Payload struct {
			Grant []struct {
				RoleAssignment RoleAssignment
			}
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	variables := map[string]any{"input": input}
	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	// The mutation returns a list of granted assignments
	if len(mutation.Payload.Grant) == 0 {
		return nil, APIError{Kind: "API Error", Detail: "No role assignment granted in response"}
	}

	return &mutation.Payload.Grant[0].RoleAssignment, nil
}

// Read returns data for a specific role assignment by ID.
func (r roleAssignmentAPI) Read(ctx context.Context, id string) (*RoleAssignment, error) {
	var query struct {
		RoleAssignments struct {
			Edges []struct {
				Node RoleAssignment
			}
		} `graphql:"roleAssignments(filterElement: $filterElement)"`
	}
	variables := map[string]any{
		"filterElement": newExactMatchFilter("id", id),
	}
	if err := r.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(query.RoleAssignments.Edges) == 0 {
		return nil, NotFound{"Role assignment not found"}
	}

	return &query.RoleAssignments.Edges[0].Node, nil
}

// List returns role assignments filtered by target or principal.
func (r roleAssignmentAPI) List(ctx context.Context, target *RoleAssignmentTarget, principal *RoleAssignmentPrincipal) ([]RoleAssignment, error) {
	var query struct {
		RoleAssignments struct {
			Edges []struct {
				Node RoleAssignment
			}
		} `graphql:"roleAssignments(filterElement: $filterElement)"`
	}

	// Build filter based on target and/or principal
	var filterElement FilterElementInput
	if target != nil {
		// For now, filter by target type - may need to be enhanced for UUID filtering
		filterElement = newExactMatchFilter("target.type", target.Type)
	} else if principal != nil {
		// Filter by principal type and ID
		filterElement = newExactMatchFilter("principal.type", principal.Type)
	}

	variables := map[string]any{
		"filterElement": filterElement,
	}

	if err := r.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	assignments := make([]RoleAssignment, 0, len(query.RoleAssignments.Edges))
	for _, edge := range query.RoleAssignments.Edges {
		// Apply client-side filtering for fields not in the query filter
		assignment := edge.Node

		// Filter by target if specified
		if target != nil {
			assignmentTarget := assignment.GetTarget()
			if assignmentTarget.Type != target.Type {
				continue
			}
			// Check UUID if provided (for non-system targets)
			if target.UUID != nil && (assignmentTarget.UUID == nil || *assignmentTarget.UUID != *target.UUID) {
				continue
			}
			// For system target, ensure UUID is nil
			if target.Type == "system" && assignmentTarget.UUID != nil {
				continue
			}
		}

		// Filter by principal if specified
		if principal != nil {
			assignmentPrincipal := assignment.GetPrincipal()
			if assignmentPrincipal.Type != principal.Type || assignmentPrincipal.ID != principal.ID {
				continue
			}
		}

		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// ListByTargetString returns role assignments filtered by an opaque target string.
// The target string should be in the format "type:id" (e.g., "account-group:uuid", "system:all").
func (r roleAssignmentAPI) ListByTargetString(ctx context.Context, targetStr string) ([]RoleAssignment, error) {
	// Get all role assignments (no server-side filtering for now)
	// We'll filter client-side by comparing the target strings
	assignments, err := r.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	// Filter by target string
	filtered := make([]RoleAssignment, 0)
	for _, assignment := range assignments {
		assignmentTarget := assignment.GetTarget()
		if assignmentTarget.String() == targetStr {
			filtered = append(filtered, assignment)
		}
	}

	return filtered, nil
}

// Delete removes a role assignment.
func (r roleAssignmentAPI) Delete(ctx context.Context, roleName string, principal RoleAssignmentPrincipal, target RoleAssignmentTarget) error {
	// Validate target
	targetStr := target.String()
	if targetStr == "" {
		return APIError{
			Kind:   "Invalid Input",
			Detail: fmt.Sprintf("Target type '%s' requires a UUID, but none was provided", target.Type),
		}
	}

	input := RoleAssignmentInput{
		Revokes: []RoleAssignmentRevoke{
			{
				Role:      roleName,
				Principal: principal.String(),
				Target:    targetStr,
			},
		},
	}

	var mutation struct {
		Payload struct {
			Revoke []struct {
				RoleAssignment RoleAssignment
			}
		} `graphql:"updateRoleAssignment(input: $input)"`
	}

	variables := map[string]any{"input": input}
	if err := r.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}

	return nil
}
