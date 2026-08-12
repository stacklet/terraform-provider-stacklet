// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// UserGroup is the data returned by reading user group data.
type UserGroup struct {
	ID                      graphql.ID `graphql:"id"`
	UUID                    string     `graphql:"uuid"`
	Name                    string     `graphql:"name"`
	DisplayName             *string    `graphql:"displayName"`
	RoleAssignmentPrincipal string     `graphql:"roleAssignmentPrincipal"`
	RoleAssignmentTarget    string     `graphql:"roleAssignmentTarget"`
}

// UserGroupCreateInput is the input for creating a user group.
type UserGroupCreateInput struct {
	Name        string  `json:"name"`
	DisplayName *string `json:"displayName,omitempty"`
}

func (i UserGroupCreateInput) GetGraphQLType() string {
	return "AddUserGroupInput"
}

// UserGroupUpdateInput is the input for updating a user group.
type UserGroupUpdateInput struct {
	ID          string  `json:"id"`
	Name        *string `json:"name"`
	DisplayName *string `json:"displayName"`
}

func (i UserGroupUpdateInput) GetGraphQLType() string {
	return "UpdateUserGroupInput"
}

// UserGroupDeleteInput is the input for removing a user group.
type UserGroupDeleteInput struct {
	ID string `json:"id"`
}

func (i UserGroupDeleteInput) GetGraphQLType() string {
	return "RemoveUserGroupInput"
}

type userGroupAPI struct {
	c *client
}

// Read returns data for a user group by its UUID.
func (a userGroupAPI) Read(ctx context.Context, uuid string) (*UserGroup, error) {
	var query struct {
		UserGroup *UserGroup `graphql:"userGroup(uuid: $uuid)"`
	}
	variables := map[string]any{
		"uuid": UUID(uuid),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}
	if query.UserGroup == nil {
		return nil, NotFound{"User group not found"}
	}
	return query.UserGroup, nil
}

// ReadByName returns data for a user group by its name.
func (a userGroupAPI) ReadByName(ctx context.Context, name string) (*UserGroup, error) {
	var query struct {
		UserGroups struct {
			Edges []struct {
				Node UserGroup
			}
		} `graphql:"userGroups(filterElement: $filterElement)"`
	}
	variables := map[string]any{
		"filterElement": newExactMatchFilter("name", name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}
	if len(query.UserGroups.Edges) == 0 {
		return nil, NotFound{"User group not found"}
	}
	return &query.UserGroups.Edges[0].Node, nil
}

// Create creates a user group.
func (a userGroupAPI) Create(ctx context.Context, i UserGroupCreateInput) (*UserGroup, error) {
	var mutation struct {
		Payload struct {
			UserGroup *UserGroup `graphql:"userGroup"`
		} `graphql:"addUserGroup(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": i}); err != nil {
		return nil, err
	}
	if mutation.Payload.UserGroup == nil {
		return nil, NotFound{"User group not found after creation"}
	}
	return mutation.Payload.UserGroup, nil
}

// Update updates a user group.
func (a userGroupAPI) Update(ctx context.Context, i UserGroupUpdateInput) (*UserGroup, error) {
	var mutation struct {
		Payload struct {
			UserGroup *UserGroup `graphql:"userGroup"`
		} `graphql:"updateUserGroup(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": i}); err != nil {
		return nil, err
	}
	if mutation.Payload.UserGroup == nil {
		return nil, NotFound{"User group not found after update"}
	}
	return mutation.Payload.UserGroup, nil
}

// Delete removes a user group.
func (a userGroupAPI) Delete(ctx context.Context, id string) error {
	var mutation struct {
		Payload struct {
			Removed struct {
				ID graphql.ID `graphql:"id"`
			} `graphql:"removed"`
		} `graphql:"removeUserGroup(input: $input)"`
	}
	input := UserGroupDeleteInput{ID: id}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": input}); err != nil {
		return err
	}
	return nil
}
