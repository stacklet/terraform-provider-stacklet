// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"
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
	BotEndpoint  string `graphql:"botEndpoint"`
	WIFIssuerURL string `graphql:"wifIssuerURL"`
	TrustRoleARN string `graphql:"trustRoleARN"`
}

// GCPIntegrationSurface is the data returned by reading GCP integration configuration details.
type GCPIntegrationSurface struct {
	TrustAws GCPIntegrationSurfaceTrustAws `graphql:"trustAws"`
	AwsRelay GCPIntegrationSurfaceAwsRelay `graphql:"awsRelay"`
}

// GCPIntegrationSurfaceTrustAws holds AWS trust configuration for GCP integration.
type GCPIntegrationSurfaceTrustAws struct {
	AccountID         string `graphql:"accountId"`
	AssetdbRoleName   string `graphql:"assetdbRoleName"`
	CostQueryRoleName string `graphql:"costQueryRoleName"`
	ExecutionRoleName string `graphql:"executionRoleName"`
	PlatformRoleName  string `graphql:"platformRoleName"`
}

// GCPIntegrationSurfaceAwsRelay holds AWS relay configuration for GCP integration.
type GCPIntegrationSurfaceAwsRelay struct {
	BusArn  string `graphql:"busArn"`
	RoleArn string `graphql:"roleArn"`
}

type systemAPI struct {
	c *client
}

// Platform returns platform details.
func (a systemAPI) Platform(ctx context.Context) (*Platform, error) {
	var query struct {
		Platform Platform `graphql:"platform"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}
	return &query.Platform, nil
}

// MSTeamsIntegrationSurface returns details for the MSTeams platform integration.
func (a systemAPI) MSTeamsIntegrationSurface(ctx context.Context) (*MSTeamsIntegrationSurface, error) {
	var query struct {
		MSTeamsIntegrationSurface MSTeamsIntegrationSurface `graphql:"msTeamsIntegrationSurface"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}
	return &query.MSTeamsIntegrationSurface, nil
}

// GCPIntegrationSurface returns details for the GCP platform integration.
func (a systemAPI) GCPIntegrationSurface(ctx context.Context) (*GCPIntegrationSurface, error) {
	var query struct {
		GCPIntegrationSurface GCPIntegrationSurface `graphql:"gcpIntegrationSurface"`
	}
	if err := a.c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}
	return &query.GCPIntegrationSurface, nil
}
