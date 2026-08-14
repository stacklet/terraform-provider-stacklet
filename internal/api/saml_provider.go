// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// SAMLProvider is the data returned by reading SAML provider data.
type SAMLProvider struct {
	ID            graphql.ID `graphql:"id"`
	Name          string     `graphql:"name"`
	DisplayName   *string    `graphql:"displayName"`
	MetadataURL   *string    `graphql:"metadataURL"`
	MetadataXML   *string    `graphql:"metadataXML"`
	EnableSignout bool       `graphql:"enableSignout"`
	IDPAlias      *string    `graphql:"idpAlias"`
}

// SAMLProviderCreateInput is the input for creating a SAML provider.
type SAMLProviderCreateInput struct {
	Name          string  `json:"name"`
	DisplayName   *string `json:"displayName,omitempty"`
	MetadataURL   *string `json:"metadataURL,omitempty"`
	MetadataXML   *string `json:"metadataXML,omitempty"`
	EnableSignout bool    `json:"enableSignout"`
	IDPAlias      *string `json:"idpAlias,omitempty"`
}

func (i SAMLProviderCreateInput) GetGraphQLType() string {
	return "AddSAMLProviderInput"
}

// SAMLProviderUpdateInput is the input for updating a SAML provider.
//
// The provider is identified by its (immutable) name. Every optional field is
// always sent, so that a null clears the corresponding value rather than
// leaving it unchanged, as the API does for omitted fields.
type SAMLProviderUpdateInput struct {
	Name          string  `json:"name"`
	DisplayName   *string `json:"displayName"`
	MetadataURL   *string `json:"metadataURL"`
	MetadataXML   *string `json:"metadataXML"`
	EnableSignout bool    `json:"enableSignout"`
	IDPAlias      *string `json:"idpAlias"`
}

func (i SAMLProviderUpdateInput) GetGraphQLType() string {
	return "UpdateSAMLProviderInput"
}

// samlProviderDeleteInput is the input for removing a SAML provider.
type samlProviderDeleteInput struct {
	Name string `json:"name"`
}

func (i samlProviderDeleteInput) GetGraphQLType() string {
	return "RemoveSAMLProviderInput"
}

type samlProviderAPI struct {
	c *client
}

// Read returns data for a SAML provider by its (unique) name.
func (a samlProviderAPI) Read(ctx context.Context, name string) (*SAMLProvider, error) {
	var query struct {
		SAMLProvider *SAMLProvider `graphql:"samlProvider(name: $name)"`
	}
	variables := map[string]any{
		"name": graphql.String(name),
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}
	if query.SAMLProvider == nil {
		return nil, NotFound{"SAML provider not found"}
	}
	return query.SAMLProvider, nil
}

// Create adds a SAML provider, which the platform reconciles into the identity pool.
func (a samlProviderAPI) Create(ctx context.Context, i SAMLProviderCreateInput) (*SAMLProvider, error) {
	var mutation struct {
		Payload struct {
			Provider *SAMLProvider `graphql:"provider"`
			Problems []problem
		} `graphql:"addSAMLProvider(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": i}); err != nil {
		return nil, err
	}
	if err := fromProblems(ctx, mutation.Payload.Problems); err != nil {
		return nil, err
	}
	if mutation.Payload.Provider == nil {
		return nil, NotFound{"SAML provider not found after creation"}
	}
	return mutation.Payload.Provider, nil
}

// Update applies changes to the SAML provider identified by name; names are immutable.
func (a samlProviderAPI) Update(ctx context.Context, i SAMLProviderUpdateInput) (*SAMLProvider, error) {
	var mutation struct {
		Payload struct {
			Provider *SAMLProvider `graphql:"provider"`
			Problems []problem
		} `graphql:"updateSAMLProvider(input: $input)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": i}); err != nil {
		return nil, err
	}
	if err := fromProblems(ctx, mutation.Payload.Problems); err != nil {
		return nil, err
	}
	if mutation.Payload.Provider == nil {
		return nil, NotFound{"SAML provider not found after update"}
	}
	return mutation.Payload.Provider, nil
}

// Delete removes a SAML provider by its name.
func (a samlProviderAPI) Delete(ctx context.Context, name string) error {
	var mutation struct {
		Payload struct {
			Problems []problem
		} `graphql:"removeSAMLProvider(input: $input)"`
	}
	input := samlProviderDeleteInput{Name: name}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"input": input}); err != nil {
		return err
	}
	return fromProblems(ctx, mutation.Payload.Problems)
}
