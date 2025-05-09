// package api provides access to the GraphQL API.
package api

import (
	"github.com/hasura/go-graphql-client"
)

// API provides access to the GraphQL API.
type API struct {
	Account accountAPI
	Policy  policyAPI
}

// New creates an API wrapper.
func New(c *graphql.Client) *API {
	return &API{
		Account: accountAPI{c},
		Policy:  policyAPI{c},
	}
}
