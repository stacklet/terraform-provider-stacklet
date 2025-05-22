package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// Account is the data returned by reading account data.
type AccountDiscovery struct {
	ID          string
	Name        string
	Description *string
	Provider    CloudProvider
	Config      struct {
		TypeName    string                      `graphql:"__typename"`
		AWSConfig   accountDiscoveryAWSConfig   `graphql:"... on AWSAccountDiscoveryConfig"`
		AzureConfig accountDiscoveryAzureConfig `graphql:"... on AzureAccountDiscoveryConfig"`
		GCPConfig   accountDiscoveryGCPConfig   `graphql:"... on GCPAccountDiscoveryConfig"`
	}
	Schedule struct {
		Suspended bool
	}
}

type accountDiscoveryAWSConfig struct {
	OrgID         string `graphql:"orgID"`
	OrgRole       string `graphql:"orgRole"`
	MemberRole    string `graphql:"memberRole"`
	CustodianRole string `graphql:"custodianRole"`
}

type accountDiscoveryAzureConfig struct {
	TenantID string `graphql:"tenantID"`
	ClientID string `graphql:"clientID"`
}

type accountDiscoveryGCPConfig struct {
	ClientEmail      string   `graphql:"clientEmail"`
	ClientID         string   `graphql:"clientID"`
	OrgID            string   `graphql:"orgID"`
	RootFolderIDs    []string `graphql:"rootFolderIDs"`
	ExcludeFolderIDs []string `graphql:"excludeFolderIDs"`
	ProjectID        string   `graphql:"projectID"`
	PrivateKeyID     string   `graphql:"privateKeyID"`
}

// AccountDiscoveryAWSInput is the input to create or update an AWS account discovery.
type AccountDiscoveryAWSInput struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	CustodianRole *string `json:"custodianRole,omitempty"`
	MemberRole    *string `json:"memberRole,omitempty"`
	OrgReadRole   string  `json:"orgReadRole,omitempty"`
}

func (i AccountDiscoveryAWSInput) GetGraphQLType() string {
	return "UpsertAWSAccountDiscoveryInput"
}

// AccountDiscoveryAzureInput is the input to create or update an Azure account discovery.
type AccountDiscoveryAzureInput struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	ClientID     string  `json:"clientID"`
	ClientSecret *string `json:"clientSecret,omitempty"`
	TenantID     string  `json:"tenantID"`
}

func (i AccountDiscoveryAzureInput) GetGraphQLType() string {
	return "UpsertAzureAccountDiscoveryInput"
}

// AccountDiscoveryGCPInput is the input to create or update a GCP account discovery.
type AccountDiscoveryGCPInput struct {
	Name             string   `json:"name"`
	Description      *string  `json:"description,omitempty"`
	OrgID            string   `json:"orgID"`
	RootFolderIDs    []string `json:"rootFolderIDs,omitempty"`
	ExcludeFolderIDs []string `json:"excludeFolderIDs,omitempty"`
	CredentialJSON   *string  `json:"credentialJSON,omitempty"`
}

func (i AccountDiscoveryGCPInput) GetGraphQLType() string {
	return "UpsertGCPAccountDiscoveryInput"
}

type updateAccountDiscoveryScheduleInput struct {
	Schedules []accountDiscoveryScheduleInput `json:"schedules"`
}

func (i updateAccountDiscoveryScheduleInput) GetGraphQLType() string {
	return "UpdateAccountDiscoveryScheduleInput"
}

type accountDiscoveryScheduleInput struct {
	Discovery graphql.ID      `json:"discovery"`
	Suspended graphql.Boolean `json:"suspended"`
}

type accountDiscoveryAPI struct {
	c *graphql.Client
}

// Read returns data for an account discovery.
func (a accountDiscoveryAPI) Read(ctx context.Context, name string) (AccountDiscovery, error) {
	var query struct {
		AccountDiscovery AccountDiscovery `graphql:"accountDiscovery(name: $name)"`
	}
	variables := map[string]any{"name": graphql.String(name)}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return query.AccountDiscovery, APIError{"Client error", err.Error()}
	}

	return query.AccountDiscovery, nil
}

// UpsertAWS creates or updates an AWS account discovery.
func (a accountDiscoveryAPI) UpsertAWS(ctx context.Context, input AccountDiscoveryAWSInput) (AccountDiscovery, error) {
	var mutation struct {
		UpsertAWSAccountDiscovery struct {
			AccountDiscovery AccountDiscovery
		} `graphql:"upsertAWSAccountDiscovery(input: $input)"`
	}
	variables := map[string]any{"input": input}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return mutation.UpsertAWSAccountDiscovery.AccountDiscovery, APIError{"Client error", err.Error()}
	}
	return mutation.UpsertAWSAccountDiscovery.AccountDiscovery, nil
}

// UpsertAzure creates or updates an Azure account discovery.
func (a accountDiscoveryAPI) UpsertAzure(ctx context.Context, input AccountDiscoveryAzureInput) (AccountDiscovery, error) {
	var mutation struct {
		UpsertAzureAccountDiscovery struct {
			AccountDiscovery AccountDiscovery
		} `graphql:"upsertAzureAccountDiscovery(input: $input)"`
	}
	variables := map[string]any{"input": input}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return mutation.UpsertAzureAccountDiscovery.AccountDiscovery, APIError{"Client error", err.Error()}
	}
	return mutation.UpsertAzureAccountDiscovery.AccountDiscovery, nil
}

// UpsertGCP creates or updates a GCP account discovery.
func (a accountDiscoveryAPI) UpsertGCP(ctx context.Context, input AccountDiscoveryGCPInput) (AccountDiscovery, error) {
	var mutation struct {
		UpsertGCPAccountDiscovery struct {
			AccountDiscovery AccountDiscovery
		} `graphql:"upsertGCPAccountDiscovery(input: $input)"`
	}
	variables := map[string]any{"input": input}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return mutation.UpsertGCPAccountDiscovery.AccountDiscovery, APIError{"Client error", err.Error()}
	}
	return mutation.UpsertGCPAccountDiscovery.AccountDiscovery, nil
}

// UpdateSuspended updates the susended flag for an account discovery.
func (a accountDiscoveryAPI) UpdateSuspended(ctx context.Context, id string, suspended bool) (AccountDiscovery, error) {
	var mutation struct {
		UpdateAccountDiscoverySchedule struct {
			AccountDiscoveries []AccountDiscovery
		} `graphql:"updateAccountDiscoverySchedule(input: $input)"`
	}
	variables := map[string]any{
		"input": updateAccountDiscoveryScheduleInput{
			Schedules: []accountDiscoveryScheduleInput{
				{
					Discovery: graphql.ID(id),
					Suspended: graphql.Boolean(suspended),
				},
			},
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return AccountDiscovery{}, APIError{"Client error", err.Error()}
	}

	return mutation.UpdateAccountDiscoverySchedule.AccountDiscoveries[0], nil
}
