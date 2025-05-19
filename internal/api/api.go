// package api provides access to the GraphQL API.
package api

import (
	"github.com/hasura/go-graphql-client"
)

// API provides access to the GraphQL API.
type API struct {
	Account             accountAPI
	AccountGroup        accountGroupAPI
	AccountGroupMapping accountGroupMappingAPI
	Policy              policyAPI
	Repository          repositoryAPI
}

// New creates an API wrapper.
func New(c *graphql.Client) *API {
	return &API{
		Account:             accountAPI{c},
		AccountGroup:        accountGroupAPI{c},
		AccountGroupMapping: accountGroupMappingAPI{c},
		Policy:              policyAPI{c},
		Repository:          repositoryAPI{c},
	}
}
