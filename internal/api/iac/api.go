// Copyright Stacklet, Inc. 2025, 2026

package iac

// API provides access to the IaC REST API.
type API struct {
	Policy           policyAPI
	PolicyCollection policyCollectionAPI
	PolicySource     policySourceAPI
	Project          projectAPI
}

// New creates an API wrapper.
func New(config ClientConfig) *API {
	c := newClient(config)
	return &API{
		Policy:           policyAPI{c},
		PolicyCollection: policyCollectionAPI{c},
		PolicySource:     policySourceAPI{c},
		Project:          projectAPI{c},
	}
}
