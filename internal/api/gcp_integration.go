// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// GCPIntegration is the data for a GCP integration.
type GCPIntegration struct {
	ID             graphql.ID
	Key            string
	CustomerConfig GCPIntegrationCustomerConfig
	AccessConfig   *GCPIntegrationAccessConfig
}

// GCPIntegrationCustomerConfig defines the customer-provided configuration for
// the GCP integration.
type GCPIntegrationCustomerConfig struct {
	Infrastructure   *GCPIntegrationCustomerInfrastructure
	Organizations    []GCPIntegrationCustomerOrganization
	CostSources      []GCPIntegrationCustomerCostSource
	SecurityContexts []GCPIntegrationCustomerSecurityContext
	TerraformModule  *TerraformModule
}

// GCPIntegrationCustomerInfrastructure defines the project configuration for
// the GCP integration deployment.
type GCPIntegrationCustomerInfrastructure struct {
	ProjectID        string `graphql:"projectId"`
	ResourceLocation string
	ResourcePrefix   string
	CreateProject    *GCPIntegrationCustomerCreateProject
}

// GCPIntegrationCustomerCreateProject holds configuration for integration
// project creation.
type GCPIntegrationCustomerCreateProject struct {
	BillingAccountID string  `graphql:"billingAccountId"`
	OrgID            *string `graphql:"orgId"`
	FolderID         *string `graphql:"folderId"`
	Labels           []Tag
}

// GCPIntegrationCustomerOrganization identifies a GCP organization Stacklet will manage.
type GCPIntegrationCustomerOrganization struct {
	OrgID      string   `graphql:"orgId"`
	FolderIDs  []string `graphql:"folderIds"`
	ProjectIDs []string `graphql:"projectIds"`
}

// GCPIntegrationCustomerCostSource identifies a billing export table for cost data.
type GCPIntegrationCustomerCostSource struct {
	BillingTable string
}

// GCPIntegrationCustomerSecurityContext defines additional security context to
// define in the integration.
type GCPIntegrationCustomerSecurityContext struct {
	Name       string
	ExtraRoles []string
}

// GCPIntegrationAccessConfig holds the access details from the GCP
// integration.
type GCPIntegrationAccessConfig struct {
	Infrastructure   GCPIntegrationAccessInfrastructure
	Organizations    []GCPIntegrationAccessOrganization
	CostSources      []GCPIntegrationAccessCostSource
	SecurityContexts []GCPIntegrationAccessSecurityContext
	RoundtripDigest  string
}

// GCPIntegrationAccessInfrastructure identifies the GCP infrastructure project.
type GCPIntegrationAccessInfrastructure struct {
	ProjectID     string `graphql:"projectId"`
	Relay         GCPIntegrationAccessRelay
	WIF           GCPIntegrationAccessWIF `graphql:"wif"`
	BaselineRoles []string
}

// GCPIntegrationAccessRelay holds the relay identity credential.
type GCPIntegrationAccessRelay struct {
	OAuthID string `graphql:"oauthId"`
}

// GCPIntegrationAccessWIF holds the Workload Identity Federation configuration.
type GCPIntegrationAccessWIF struct {
	Audience   string
	Principals GCPIntegrationAccessWIFPrincipals
}

// GCPIntegrationAccessWIFPrincipals identifies WIF principals by their intended role.
type GCPIntegrationAccessWIFPrincipals struct {
	ReadOnly  string
	CostQuery string
}

// GCPIntegrationAccessOrganization identifies an accessible GCP organization.
type GCPIntegrationAccessOrganization struct {
	ID       string
	Name     string
	Folders  []GCPIntegrationAccessOrganizationFolder
	Projects []GCPIntegrationAccessOrganizationProject
}

// GCPIntegrationAccessOrganizationFolder provides details about a connected organization folder.
type GCPIntegrationAccessOrganizationFolder struct {
	ID   string
	Name string
}

// GCPIntegrationAccessOrganizationProject provides details about a connected organization project.
type GCPIntegrationAccessOrganizationProject struct {
	ID     string
	Number string
}

// GCPIntegrationAccessCostSource identifies a billing export table and its location.
type GCPIntegrationAccessCostSource struct {
	BillingTable string
	Location     string
}

// GCPIntegrationAccessSecurityContext defines a named set of roles granted to a principal.
type GCPIntegrationAccessSecurityContext struct {
	Name       string
	ExtraRoles []string
	Principal  string
}

type gcpIntegrationAPI struct {
	c *graphql.Client
}

// Read returns data for a GCP integration by key.
func (a gcpIntegrationAPI) Read(ctx context.Context, key string) (*GCPIntegration, error) {
	var query struct {
		Payload struct {
			GCPIntegration *GCPIntegration `graphql:"gcpIntegration"`
			Problems       []Problem
		} `graphql:"gcpIntegration(key: $key)"`
	}
	if err := a.c.Query(ctx, &query, map[string]any{"key": graphql.String(key)}); err != nil {
		return nil, NewAPIError(err)
	}
	if err := fromProblems(ctx, query.Payload.Problems); err != nil {
		return nil, err
	}
	if query.Payload.GCPIntegration == nil {
		return nil, NotFound{"GCP integration not found"}
	}
	return query.Payload.GCPIntegration, nil
}
