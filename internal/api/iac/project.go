// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"context"
)

// ProjectGroups holds the groups with access to a project.
type ProjectGroups struct {
	Read []string `json:"read"`
}

// Project is the data returned by reading project data.
type Project struct {
	Name        string         `json:"name"`
	ClientID    *string        `json:"client_id"`
	Collections []string       `json:"collections"`
	Groups      *ProjectGroups `json:"groups"`
}

// ProjectInput is the input for creating or updating a project.
type ProjectInput struct {
	Collections []string      `json:"collections"`
	Groups      ProjectGroups `json:"groups"`
}

type projectCreateInput struct {
	Name string `json:"name"`
	ProjectInput
}

type projectAPI struct {
	c *client
}

// Read returns data for a project by name.
func (p projectAPI) Read(ctx context.Context, name string) (*Project, error) {
	var project Project
	if err := p.c.Get(ctx, resourcePath("projects", name), nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// Create creates a project.
func (p projectAPI) Create(ctx context.Context, name string, input ProjectInput) (*Project, error) {
	var project Project
	payload := projectCreateInput{Name: name, ProjectInput: input}
	if err := p.c.Post(ctx, "projects", payload, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// Update updates a project.
func (p projectAPI) Update(ctx context.Context, name string, input ProjectInput) (*Project, error) {
	var project Project
	if err := p.c.Put(ctx, resourcePath("projects", name), input, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// Delete deletes a project.
func (p projectAPI) Delete(ctx context.Context, name string) error {
	return p.c.Delete(ctx, resourcePath("projects", name), nil)
}
