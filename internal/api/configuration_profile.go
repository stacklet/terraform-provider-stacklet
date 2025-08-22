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
		SlackConfiguration      SlackConfiguration      `graphql:"... on SlackConfiguration"`
		SymphonyConfiguration   SymphonyConfiguration   `graphql:"... on SymphonyConfiguration"`
		TeamsConfiguration      TeamsConfiguration      `graphql:"... on TeamsConfiguration"`
		JiraConfiguration       JiraConfiguration       `graphql:"... on JiraConfiguration"`
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
	IssueType   string
	ClosedState string
}

// SymphonyConfiguration is the configuration for Symphony profiles.
type SymphonyConfiguration struct {
	AgentDomain    string
	ServiceAccount string
}

// SlackConfiguration is the configuration for Symphony profiles.
type SlackConfiguration struct {
	UserFields []string
	Webhooks   []SlackWebhook
}

// SlackWebhook is a webhook configuration for Slack.
type SlackWebhook struct {
	Name string
	URL  string `graphql:"url"`
}

// TeamsConfiguration is the configuration for Microsoft Teams profiles.
type TeamsConfiguration struct {
	Webhooks []TeamsWebhook
}

// TeamsWebhook is a webhook configuration for Microsoft Teams.
type TeamsWebhook struct {
	Name string
	URL  string `graphql:"url"`
}

// JiraConfiguation is the configuration for Jira profiles.
type JiraConfiguration struct {
	URL      *string `graphql:"url"`
	Projects []JiraProject
	User     string
}

// JiraPorject is the configuration for a Jira project.
type JiraProject struct {
	ClosedStatus string
	IssueType    string
	Name         string
	Project      string
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

// ReadEmail returns data for the email configuration profile.
func (a configurationProfileAPI) ReadEmail(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileEmail)
}

// ReadSlack returns data for the Slack configuration profile.
func (a configurationProfileAPI) ReadSlack(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileSlack)
}

// ReadTeams returns data for the Microsoft Teams configuration profile.
func (a configurationProfileAPI) ReadTeams(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileTeams)
}

// ReadSymphony returns data for the Symphony configuration profile.
func (a configurationProfileAPI) ReadSymphony(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileSymphony)
}

// ReadSymphony returns data for the ServiceNow configuration profile.
func (a configurationProfileAPI) ReadServiceNow(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileServiceNow)
}

// ReadJira returns data for the Jira configuration profile.
func (a configurationProfileAPI) ReadJira(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileJira)
}
