// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Template is the data returned for a notification template.
type Template struct {
	ID          string
	Name        string
	Description *string
	Transport   *string
	Content     string
}

// TemplateCreateInput is the input for creating a template.
type TemplateInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Transport   *string `json:"transport"`
	Content     string  `json:"content"`
}

type templateAPI struct {
	c *graphql.Client
}

// Read returns data for a notification template.
func (a templateAPI) Read(ctx context.Context, name string) (*Template, error) {
	var query struct {
		Template Template `graphql:"template(name: $name)"`
	}

	if err := a.c.Query(ctx, &query, map[string]any{"name": graphql.String(name)}); err != nil {
		return nil, NewAPIError(err)
	}
	if query.Template.ID == "" {
		return nil, NotFound{"Template not found"}
	}

	return &query.Template, nil
}

// Upsert creates or updates a notification template.
func (a templateAPI) Upsert(ctx context.Context, input TemplateInput) (*Template, error) {
	var mutation struct {
		Payload struct {
			Template Template
		} `graphql:"addTemplate(input: $input)"`
	}
	err := a.c.Mutate(ctx, &mutation, map[string]any{"input": input})
	if err != nil {
		return nil, NewAPIError(err)
	}

	return &mutation.Payload.Template, nil
}

// Delete removes a notification template.
func (a templateAPI) Delete(ctx context.Context, name string) error {
	var mutation struct {
		ID string `graphql:"removeTemplate(name: $name)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"name": graphql.String(name)}); err != nil {
		return NewAPIError(err)
	}
	return nil
}
