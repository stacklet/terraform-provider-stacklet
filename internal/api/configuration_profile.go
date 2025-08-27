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
		TypeName                   string                     `graphql:"__typename"`
		EmailConfiguration         EmailConfiguration         `graphql:"... on EmailConfiguration"`
		ServiceNowConfiguration    ServiceNowConfiguration    `graphql:"... on ServiceNowConfiguration"`
		SlackConfiguration         SlackConfiguration         `graphql:"... on SlackConfiguration"`
		SymphonyConfiguration      SymphonyConfiguration      `graphql:"... on SymphonyConfiguration"`
		TeamsConfiguration         TeamsConfiguration         `graphql:"... on TeamsConfiguration"`
		JiraConfiguration          JiraConfiguration          `graphql:"... on JiraConfiguration"`
		ResourceOwnerConfiguration ResourceOwnerConfiguration `graphql:"... on ResourceOwnerConfiguration"`
		AccountOwnersConfiguration AccountOwnersConfiguration `graphql:"... on AccountOwnersConfiguration"`
	}
}

// EmailConfiguration is the configuration for email profiles.
type EmailConfiguration struct {
	FromEmail string
	SESRegion *string            `graphql:"sesRegion"`
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
	Endpoint    string `json:"endpoint"`
	User        string `json:"user"`
	Password    string `json:"password"`
	IssueType   string `json:"issueType"`
	ClosedState string `json:"closedState"`
}

// SymphonyConfiguration is the configuration for Symphony profiles.
type SymphonyConfiguration struct {
	AgentDomain    string
	ServiceAccount string
}

// SlackConfiguration is the configuration for Slack profiles.
type SlackConfiguration struct {
	Token      *string        `json:"token"`
	UserFields []string       `json:"userFields"`
	Webhooks   []SlackWebhook `json:"webhooks"`
}

// SlackWebhook is a webhook configuration for Slack.
type SlackWebhook struct {
	Name string `json:"name"`
	URL  string `graphql:"url" json:"url"`
}

// TeamsConfiguration is the configuration for Microsoft Teams profiles.
type TeamsConfiguration struct {
	Webhooks []TeamsWebhook `json:"webhooks"`
}

// TeamsWebhook is a webhook configuration for Microsoft Teams.
type TeamsWebhook struct {
	Name string `json:"name"`
	URL  string `graphql:"url" json:"url"`
}

// JiraConfiguation is the configuration for Jira profiles.
type JiraConfiguration struct {
	URL      *string       `graphql:"url" json:"url"`
	Projects []JiraProject `json:"projects"`
	User     string        `json:"user"`
	APIKey   string        `json:"apiKey"`
}

// JiraPorject is the configuration for a Jira project.
type JiraProject struct {
	ClosedStatus string `json:"closedStatus"`
	IssueType    string `json:"issueType"`
	Name         string `json:"name"`
	Project      string `json:"project"`
}

type jiraConfigurationInput struct {
	JiraConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i jiraConfigurationInput) GetGraphQLType() string {
	return "JiraConfigurationInput"
}

// ResourceOwnerConfiguration is the configuation for resource owner.
type ResourceOwnerConfiguration struct {
	// "default" is present with different type in both resource and account, so it must be aliased
	Default      []string `graphql:"resourceOwnerDefault: default" json:"default"`
	OrgDomain    *string  `json:"orgDomain"`
	OrgDomainTag *string  `json:"orgDomainTag"`
	Tags         []string `json:"tags"`
}

type resourceOwnerConfigurationInput struct {
	ResourceOwnerConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i resourceOwnerConfigurationInput) GetGraphQLType() string {
	return "ResourceOwnerConfigurationInput"
}

// AccountOwnersConfiguration is the configuration for account owners.
type AccountOwnersConfiguration struct {
	// "default" is present with different type in both resource and account, so it must be aliased
	Default      []AccountOwners `graphql:"accountOwnersDefault: default" json:"default"`
	OrgDomain    *string         `json:"orgDomain"`
	OrgDomainTag *string         `json:"orgDomainTag"`
	Tags         []string        `json:"tags"`
}

// AccountOwners tracks the owners for an account.
type AccountOwners struct {
	Account string   `json:"account"`
	Owners  []string `json:"owners"`
}

type accountOwnersConfigurationInput struct {
	AccountOwnersConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i accountOwnersConfigurationInput) GetGraphQLType() string {
	return "AccountOwnersConfigurationInput"
}

type teamsConfigurationInput struct {
	TeamsConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i teamsConfigurationInput) GetGraphQLType() string {
	return "TeamsConfigurationInput"
}

type slackConfigurationInput struct {
	SlackConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i slackConfigurationInput) GetGraphQLType() string {
	return "SlackConfigurationInput"
}

type serviceNowConfigurationInput struct {
	ServiceNowConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i serviceNowConfigurationInput) GetGraphQLType() string {
	return "ServiceNowConfigurationInput"
}

type configurationProfileAPI struct {
	c *graphql.Client
}

const configurationScopeGlobal = "0"

// Read returns data for a configuration profile.
func (a configurationProfileAPI) Read(ctx context.Context, name ConfigurationProfileName) (*ConfigurationProfile, error) {
	var query struct {
		Configuration ConfigurationProfile `graphql:"profile(name: $name, scope: $scope)"`
	}
	variables := map[string]any{
		"name":  graphql.String(string(name)),
		"scope": graphql.String(configurationScopeGlobal),
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

// ReadAccountResourceOwners returns data for the account owners configuration profile.
func (a configurationProfileAPI) ReadAccountOwners(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileAccountOwners)
}

// ReadResourceOwner returns data for the resource owners configuration profile.
func (a configurationProfileAPI) ReadResourceOwner(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileResourceOwner)
}

// Upsert the Jira configuration profile.
func (a configurationProfileAPI) UpsertJira(ctx context.Context, input JiraConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addJiraProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": jiraConfigurationInput{
			JiraConfiguration: input,
			Name:              string(ConfigurationProfileJira),
			Scope:             configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// Upsert the account owners configuration profile.
func (a configurationProfileAPI) UpsertAccountOwners(ctx context.Context, input AccountOwnersConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addAccountOwnersProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": accountOwnersConfigurationInput{
			AccountOwnersConfiguration: input,
			Name:                       string(ConfigurationProfileAccountOwners),
			Scope:                      configurationScopeGlobal,
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// Upsert the resource owners configuration profile.
func (a configurationProfileAPI) UpsertResourceOwner(ctx context.Context, input ResourceOwnerConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addResourceOwnerProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": resourceOwnerConfigurationInput{
			ResourceOwnerConfiguration: input,
			Name:                       string(ConfigurationProfileResourceOwner),
			Scope:                      configurationScopeGlobal,
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// UpsertTeams upserts the Microsoft Teams configuration profile.
func (a configurationProfileAPI) UpsertTeams(ctx context.Context, config TeamsConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addTeamsProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": teamsConfigurationInput{
			TeamsConfiguration: config,
			Name:               string(ConfigurationProfileTeams),
			Scope:              configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// UpsertSlack upserts the Slack configuration profile.
func (a configurationProfileAPI) UpsertSlack(ctx context.Context, config SlackConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addSlackProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": slackConfigurationInput{
			SlackConfiguration: config,
			Name:               string(ConfigurationProfileSlack),
			Scope:              configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// UpsertServiceNow upserts the ServiceNow configuration profile.
func (a configurationProfileAPI) UpsertServiceNow(ctx context.Context, config ServiceNowConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addServiceNowProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": serviceNowConfigurationInput{
			ServiceNowConfiguration: config,
			Name:                    string(ConfigurationProfileServiceNow),
			Scope:                   configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// Delete removes a configuation profile.
func (a configurationProfileAPI) Delete(ctx context.Context, name ConfigurationProfileName) error {
	var mutation struct {
		ID string `graphql:"removeProfile(scope: $scope, name: $name)"`
	}
	variables := map[string]any{
		"name":  graphql.String(string(name)),
		"scope": graphql.String(configurationScopeGlobal),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return NewAPIError(err)
	}
	return nil
}

// DeleteJira deletes the Jira configuration profile.
func (a configurationProfileAPI) DeleteJira(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileJira)
}

// DeleteTeams deletes the Microsoft Teams configuration profile.
func (a configurationProfileAPI) DeleteTeams(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileTeams)
}

// DeleteSlack deletes the Slack configuration profile.
func (a configurationProfileAPI) DeleteSlack(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileSlack)
}

// DeleteServiceNow deletes the ServiceNow configuration profile.
func (a configurationProfileAPI) DeleteServiceNow(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileServiceNow)
}

// DeleteAccountOwners deletes the account owners configuration profile.
func (a configurationProfileAPI) DeleteAccountOwners(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileAccountOwners)
}

// DeleteResourceOwner deletes the resource owners configuration profile.
func (a configurationProfileAPI) DeleteResourceOwner(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileResourceOwner)
}
