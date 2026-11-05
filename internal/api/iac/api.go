// Copyright Stacklet, Inc. 2025, 2026

package iac

// API provides access to the IaC REST API.
type API struct {
	PolicySource policySourceAPI
}

// New creates an API wrapper.
func New(config ClientConfig) *API {
	c := newClient(config)
	return &API{
		PolicySource: policySourceAPI{c},
	}
}
