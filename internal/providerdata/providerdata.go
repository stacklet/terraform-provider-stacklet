// Copyright (c) 2025 - Stacklet, Inc.

package providerdata

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// providerData holds shared data for available in requests.
type providerData struct {
	API *api.API
}

// New returns configured provider data.
func New(client *graphql.Client) *providerData {
	return &providerData{
		API: api.New(client),
	}
}

type providerDataError struct {
	Kind         string
	ProviderData any
}

func (e providerDataError) Summary() string {
	return fmt.Sprintf("Unexpected %s configure type", e.Kind)
}

func (e providerDataError) Error() string {
	return fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", e.ProviderData)
}

// GetForResource returns provider data for a resource request, or nil if not set.
func GetForResource(req resource.ConfigureRequest) (*providerData, error) {
	if req.ProviderData == nil {
		return nil, nil
	}
	if providerData, ok := req.ProviderData.(*providerData); ok {
		return providerData, nil
	}
	return nil, providerDataError{
		Kind:         "resource",
		ProviderData: req.ProviderData,
	}
}

// GetForDataSource returns provider data for a data source request, or nil if not set.
func GetForDataSource(req datasource.ConfigureRequest) (*providerData, error) {
	if req.ProviderData == nil {
		return nil, nil
	}
	if providerData, ok := req.ProviderData.(*providerData); ok {
		return providerData, nil
	}
	return nil, providerDataError{
		Kind:         "datasource",
		ProviderData: req.ProviderData,
	}
}
