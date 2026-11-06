// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"context"
)

// PolicyCollection is a collection of policies.
type PolicyCollection struct {
	Name      string   `json:"name"`
	Policies  []string `json:"policies"`
	IsDefault bool     `json:"is_default"`
	IsMirror  bool     `json:"is_mirror"`
}

// PolicyCollectionCreateInput is the input for creating a policy collection.
type PolicyCollectionCreateInput struct {
	Name      string   `json:"name"`
	Policies  []string `json:"policies"`
	IsDefault bool     `json:"is_default"`
}

// PolicyCollectionUpdateInput is the input for updating a policy collection.
type PolicyCollectionUpdateInput struct {
	Policies  []string `json:"policies"`
	IsDefault bool     `json:"is_default"`
}

type policyCollectionAPI struct {
	c *client
}

func (a policyCollectionAPI) Read(ctx context.Context, name string) (*PolicyCollection, error) {
	var collection PolicyCollection
	if err := a.c.Get(ctx, resourcePath("policy-collections", name), nil, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (a policyCollectionAPI) Create(ctx context.Context, i PolicyCollectionCreateInput) (*PolicyCollection, error) {
	i.Policies = nonNil(i.Policies)
	var collection PolicyCollection
	if err := a.c.Post(ctx, "policy-collections", i, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (a policyCollectionAPI) Update(ctx context.Context, name string, i PolicyCollectionUpdateInput) (*PolicyCollection, error) {
	// policies are left unchanged by the API if null, so always send a list
	i.Policies = nonNil(i.Policies)
	var collection PolicyCollection
	if err := a.c.Patch(ctx, resourcePath("policy-collections", name), i, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (a policyCollectionAPI) Delete(ctx context.Context, name string) error {
	return a.c.Delete(ctx, resourcePath("policy-collections", name), nil)
}

// nonNil returns an empty slice if s is nil, so that it's encoded as an empty list.
func nonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
