// Copyright Stacklet, Inc. 2025, 2026

package iac

// API provides access to the IaC REST API.
type API struct {
	c *client
}

// New creates an API wrapper.
func New(config ClientConfig) *API {
	return &API{
		c: newClient(config),
	}
}
