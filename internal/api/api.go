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
	BindingExecutionConfig  bindingExecutionConfigAPI
	Policy                  policyAPI
	PolicyCollection        policyCollectionAPI
	PolicyCollectionMapping policyCollectionMappingAPI
	Repository              repositoryAPI
}

// New creates an API wrapper.
func New(c *graphql.Client) *API {
	return &API{
		Account:                 accountAPI{c},
		AccountDiscovery:        accountDiscoveryAPI{c},
		AccountGroup:            accountGroupAPI{c},
		AccountGroupMapping:     accountGroupMappingAPI{c},
		Binding:                 bindingAPI{c},
		BindingExecutionConfig:  bindingExecutionConfigAPI{c},
		Policy:                  policyAPI{c},
		PolicyCollection:        policyCollectionAPI{c},
		PolicyCollectionMapping: policyCollectionMappingAPI{c},
		Repository:              repositoryAPI{c},
	}
}
