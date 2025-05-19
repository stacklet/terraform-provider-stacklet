package providerdata

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// ProviderData holds shared data for available in requests.
type ProviderData struct {
	API           *api.API
	GraphQLClient *graphql.Client // XXX todo remove once all are migrated
}

// New returns a new ProviderData.
func New(client *graphql.Client) *ProviderData {
	return &ProviderData{
		GraphQLClient: client,
		API:           api.New(client),
	}
}

// ProviderDataError is an error handling provider data.
type ProviderDataError struct {
	Kind         string
	ProviderData any
}

// Summary returns the error summary.
func (e ProviderDataError) Summary() string {
	return fmt.Sprintf("Unexpected %s configure type", e.Kind)
}

// Error returns the error message.
func (e ProviderDataError) Error() string {
	return fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", e.ProviderData)
}

// GetResourceProviderData returns ProviderData for a resource request, or nil if not set.
func GetResourceProviderData(req resource.ConfigureRequest) (*ProviderData, error) {
	if req.ProviderData == nil {
		return nil, nil
	}
	if providerData, ok := req.ProviderData.(*ProviderData); ok {
		return providerData, nil
	}
	return nil, ProviderDataError{
		Kind:         "resource",
		ProviderData: req.ProviderData,
	}
}

// GetDataSourceProviderData returns ProviderData for a resource request, or nil if not set.
func GetDataSourceProviderData(req datasource.ConfigureRequest) (*ProviderData, error) {
	if req.ProviderData == nil {
		return nil, nil
	}
	if providerData, ok := req.ProviderData.(*ProviderData); ok {
		return providerData, nil
	}
	return nil, ProviderDataError{
		Kind:         "datasource",
		ProviderData: req.ProviderData,
	}
}
