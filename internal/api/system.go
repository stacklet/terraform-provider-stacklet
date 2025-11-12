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
	AWSOrgReadCustomerConfig PlatformCustomerConfig `graphql:"awsOrgReadCustomerConfig"`
	AWSAccountCustomerConfig PlatformCustomerConfig `graphql:"awsAccountCustomerConfig"`
}

// PlatformCustomerConfig is the data returned for a customer configuration.
type PlatformCustomerConfig struct {
	TerraformModule TerraformModule
}

// MSTeamsIntegrationSurface is the data returned by reading Microsoft Teams
// integration configuration details.
type MSTeamsIntegrationSurface struct {
	BotEndpoint string `graphql:"botEndpoint"`
	OIDCClient  string `graphql:"oidcClient"`
	OIDCIssuer  string `graphql:"oidcIssuer"`
}

type systemAPI struct {
	c *graphql.Client
}

// Platform return platform details.
func (a systemAPI) Platform(ctx context.Context) (*Platform, error) {
	var query struct {
		Platform Platform `graphql:"platform"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, NewAPIError(err)
	}
	return &query.Platform, nil
}

// MSTeamsIntegrationSurface returns details for the MSTeams platform integration.
func (a systemAPI) MSTeamsIntegrationSurface(ctx context.Context) (*MSTeamsIntegrationSurface, error) {
	var query struct {
		MSTeamsIntegrationSurface MSTeamsIntegrationSurface `graphql:"msTeamsIntegrationSurface"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, NewAPIError(err)
	}
	return &query.MSTeamsIntegrationSurface, nil
}
