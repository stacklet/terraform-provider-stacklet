// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

// SSOGroup is the data returned by reading SSO group data.
type SSOGroup struct {
	ID                      string
	DisplayName             *string
	Name                    string
	RoleAssignmentPrincipal string
}

// SSOGroupInput is the input to create or update an SSO group.
type SSOGroupInput struct {
	Name        string  `json:"name"`
	DisplayName *string `json:"displayName"`
}

func (i SSOGroupInput) GetGraphQLType() string {
	return "UpsertSSOGroupInput"
}

type upsertSSOGroupInput struct {
	Groups []SSOGroupInput `json:"groups"`
}

func (i upsertSSOGroupInput) GetGraphQLType() string {
	return "UpsertSSOGroupsInput"
}

type upsertSSOGroupPayload struct {
	ErrorMessage *string
	SSOGroup     *SSOGroup
}

func (p upsertSSOGroupPayload) Error() string {
	if p.ErrorMessage == nil {
		return ""
	}

	return *p.ErrorMessage
}

type removeSSOGroupPayload struct {
	ErrorMessage *string
}

func (p removeSSOGroupPayload) Error() string {
	if p.ErrorMessage == nil {
		return ""
	}

	return *p.ErrorMessage
}

type removeSSOGroupsInput struct {
	Names []string `json:"names"`
}

func (i removeSSOGroupsInput) GetGraphQLType() string {
	return "RemoveSSOGroupsInput"
}

type ssoGroupAPI struct {
	c *graphql.Client
}

// Read returns data for an SSO group by name.
// Note: The Stacklet API SSO group filter requires the field name "name" with no operator.
func (a ssoGroupAPI) Read(ctx context.Context, name string) (*SSOGroup, error) {
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
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(query.SSOGroups.Edges) == 0 {
		return nil, NotFound{"SSO group not found"}
	}

	return &query.SSOGroups.Edges[0].Node, nil
}

// Upsert creates or updates an SSO group.
func (a ssoGroupAPI) Upsert(ctx context.Context, input SSOGroupInput) (*SSOGroup, error) {
	var mutation struct {
		Payload struct {
			Response []upsertSSOGroupPayload
		} `graphql:"upsertSSOGroups(input: $input)"`
	}
	variables := map[string]any{
		"input": upsertSSOGroupInput{
			Groups: []SSOGroupInput{input},
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(mutation.Payload.Response) == 0 {
		return nil, NotFound{"SSO group not found after upsert"}
	}
	payload := mutation.Payload.Response[0]
	if payload.Error() != "" {
		return nil, NewAPIError(fmt.Errorf("failed to upsert SSO group: %w", payload))
	}
	if payload.SSOGroup == nil {
		return nil, NotFound{"SSO group not found after upsert"}
	}
	return payload.SSOGroup, nil
}

// Delete removes an SSO group.
func (a ssoGroupAPI) Delete(ctx context.Context, name string) error {
	var mutation struct {
		Payload struct {
			Response []removeSSOGroupPayload
		} `graphql:"removeSSOGroups(input: $input)"`
	}
	variables := map[string]any{
		"input": removeSSOGroupsInput{Names: []string{name}},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	if len(mutation.Payload.Response) > 0 {
		payload := mutation.Payload.Response[0]
		if payload.Error() != "" {
			return NewAPIError(fmt.Errorf("failed to remove SSO group: %w", payload))
		}
	}
	return nil
}
