// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// ConfigurationProfile is the data returned for configuration profiles.
type ConfigurationProfile struct {
	ID      string
	Profile string
	Record  struct {
		TypeName                string                  `graphql:"__typename"`
		EmailConfiguration      EmailConfiguration      `graphql:"... on EmailConfiguration"`
		ServiceNowConfiguration ServiceNowConfiguration `graphql:"... on ServiceNowConfiguration"`
		SymphonyConfiguration   SymphonyConfiguration   `graphql:"... on SymphonyConfiguration"`
	}
}

// EmailConfiguration is the configuration for email profiles.
type EmailConfiguration struct {
	FromEmail string
	SesRegion *string
	SMTP      *SMTPConfiguration `graphql:"smtp"`
}

// SMTPConfiguration is the SMTP server configuration.
type SMTPConfiguration struct {
	Server   string
	Port     string
	SSL      *bool `graphql:"ssl"`
	Username *string
}

// ServiceNowConfiguration is the configuration for ServiceNow profiles.
type ServiceNowConfiguration struct {
	Endpoint    string
	User        string
	Password    string
	IssueType   string
	ClosedState string
}

// SymphonyConfiguration is the configuration for Symphony profiles.
type SymphonyConfiguration struct {
	AgentDomain    string
	ServiceAccount string
}

type configurationProfileAPI struct {
	c *graphql.Client
}

// Read returns data for a configuration profile.
func (a configurationProfileAPI) Read(ctx context.Context, name ConfigurationProfileName) (*ConfigurationProfile, error) {
	var query struct {
		Configuration ConfigurationProfile `graphql:"profile(name: $name, scope: $scope)"`
	}
	variables := map[string]any{
		"name":  graphql.String(string(name)),
		"scope": graphql.String("0"), // always use the global scope
	}
	if err := a.c.Query(ctx, &query, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if query.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found"}
	}

	return &query.Configuration, nil
}
