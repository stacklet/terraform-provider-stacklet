// Copyright (c) 2025 - Stacklet, Inc.

// package api provides access to the GraphQL API.
package api

import (
	"github.com/hasura/go-graphql-client"
)

// API provides access to the GraphQL API.
type API struct {
	Account                 accountAPI
	AccountDiscovery        accountDiscoveryAPI
	AccountGroup            accountGroupAPI
	AccountGroupMapping     accountGroupMappingAPI
	Binding                 bindingAPI
	ConfigurationProfile    configurationProfileAPI
	Platform                platformAPI
	Policy                  policyAPI
	PolicyCollection        policyCollectionAPI
	PolicyCollectionMapping policyCollectionMappingAPI
	ReportGroup             reportGroupAPI
	Repository              repositoryAPI
	Template                templateAPI
}

// New creates an API wrapper.
func New(c *graphql.Client) *API {
	return &API{
		Account:                 accountAPI{c},
		AccountDiscovery:        accountDiscoveryAPI{c},
		AccountGroup:            accountGroupAPI{c},
		AccountGroupMapping:     accountGroupMappingAPI{c},
		Binding:                 bindingAPI{c},
		ConfigurationProfile:    configurationProfileAPI{c},
		Platform:                platformAPI{c},
		Policy:                  policyAPI{c},
		PolicyCollection:        policyCollectionAPI{c},
		PolicyCollectionMapping: policyCollectionMappingAPI{c},
		ReportGroup:             reportGroupAPI{c},
		Repository:              repositoryAPI{c},
		Template:                templateAPI{c},
	}
}
