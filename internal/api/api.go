// package api provides access to the GraphQL API.
package api

import (
	"github.com/hasura/go-graphql-client"
)

// API provides access to the GraphQL API.
type API struct {
	Account      accountAPI
	AccountGroup accountGroupAPI
	Policy       policyAPI
}

// New creates an API wrapper.
func New(c *graphql.Client) *API {
	return &API{
		Account:      accountAPI{c},
		AccountGroup: accountGroupAPI{c},
		Policy:       policyAPI{c},
	}
}
