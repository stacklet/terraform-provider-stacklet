// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Platform is the data returned by reading platform data.
type Platform struct {
	ID                       string
	ExternalID               *string `graphql:"externalID"`
	ExecutionRegions         []string
	DefaultRole              *string
	AWSOrgReadCustomerConfig PlatformCustomerConfig `graphql:"awsOrgReadCustomerConfig"`
	AWSAccountCustomerConfig PlatformCustomerConfig `graphql:"awsAccountCustomerConfig"`
}

// PlatformCustomerConfig is the data returned for a customer configuration.
type PlatformCustomerConfig struct {
	TerraformModule TerraformModule
}

// TerraformModule is the data returned for terraform module definitions.
type TerraformModule struct {
	RepositoryURL string `graphql:"repositoryURL"`
	Source        string
	VariablesJSON string `graphql:"variablesJSON"`
}

type platformAPI struct {
	c *graphql.Client
}

// Read returns platform data.
func (a platformAPI) Read(ctx context.Context) (*Platform, error) {
	var query struct {
		Platform Platform `graphql:"platform"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, NewAPIError(err)
	}
	return &query.Platform, nil
}
