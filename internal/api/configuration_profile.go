// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// UUID represents a UUID scalar type in GraphQL.
type UUID string

// ConfigurationProfile is the data returned for configuration profiles.
type ConfigurationProfile struct {
	ID      graphql.ID `graphql:"id"`
	Profile string     `graphql:"profile"`
	Record  struct {
		TypeName                   string                     `graphql:"__typename"`
		EmailConfiguration         EmailConfiguration         `graphql:"... on EmailConfiguration"`
		ServiceNowConfiguration    ServiceNowConfiguration    `graphql:"... on ServiceNowConfiguration"`
		SlackConfiguration         SlackConfiguration         `graphql:"... on SlackConfiguration"`
		SymphonyConfiguration      SymphonyConfiguration      `graphql:"... on SymphonyConfiguration"`
		MSTeamsConfiguration       MSTeamsConfiguration       `graphql:"... on MSTeamsConfiguration"`
		JiraConfiguration          JiraConfiguration          `graphql:"... on JiraConfiguration"`
		ResourceOwnerConfiguration ResourceOwnerConfiguration `graphql:"... on ResourceOwnerConfiguration"`
		AccountOwnersConfiguration AccountOwnersConfiguration `graphql:"... on AccountOwnersConfiguration"`
	}
}

// EmailConfiguration is the configuration for email profiles.
type EmailConfiguration struct {
	FromEmail string             `graphql:"fromEmail" json:"fromEmail"`
	SESRegion *string            `graphql:"sesRegion" json:"sesRegion"`
	SMTP      *SMTPConfiguration `graphql:"smtp" json:"smtp"`
}

// SMTPConfiguration is the SMTP server configuration.
type SMTPConfiguration struct {
	Server   string  `graphql:"server" json:"server"`
	Port     string  `graphql:"port" json:"port"`
	SSL      *bool   `graphql:"ssl" json:"ssl"`
	Username *string `graphql:"username" json:"username"`
	Password *string `graphql:"password" json:"password"`
}

// ServiceNowConfiguration is the configuration for ServiceNow profiles.
type ServiceNowConfiguration struct {
	Endpoint    string `graphql:"endpoint" json:"endpoint"`
	User        string `graphql:"user" json:"user"`
	Password    string `graphql:"password" json:"password"`
	IssueType   string `graphql:"issueType" json:"issueType"`
	ClosedState string `graphql:"closedState" json:"closedState"`
}

// SymphonyConfiguration is the configuration for Symphony profiles.
type SymphonyConfiguration struct {
	AgentDomain    string `graphql:"agentDomain" json:"agentDomain"`
	ServiceAccount string `graphql:"serviceAccount" json:"serviceAccount"`
	PrivateKey     string `graphql:"privateKey" json:"privateKey"`
}

// SlackConfiguration is the configuration for Slack profiles.
type SlackConfiguration struct {
	Token      *string        `graphql:"token" json:"token"`
	UserFields []string       `graphql:"userFields" json:"userFields"`
	Webhooks   []SlackWebhook `graphql:"webhooks" json:"webhooks"`
}

// SlackWebhook is a webhook configuration for Slack.
type SlackWebhook struct {
	Name string `graphql:"name" json:"name"`
	URL  string `graphql:"url" json:"url"`
}

// MSTeamsConfiguration is the configuration for Microsoft Teams profiles.
type MSTeamsConfiguration struct {
	AccessConfig    *MSTeamsAccessConfig    `graphql:"accessConfig"`
	CustomerConfig  MSTeamsCustomerConfig   `graphql:"customerConfig"`
	ChannelMappings []MSTeamsChannelMapping `graphql:"channelMappings"`
	EntityDetails   MSTeamsEntityDetails    `graphql:"entityDetails"`
}

// MSTeamsAccessConfigInput is the input for the access configuration for
// Microsoft Teams profile setup.
type MSTeamsAccessConfigInput struct {
	ClientID        string `json:"clientId"`
	RoundtripDigest string `json:"roundtripDigest"`
	TenantID        string `json:"tenantId"`
}

// MSTeamsCustomerConfigInput is the input for the Microsoft teams customer configuration.
type MSTeamsCustomerConfigInput struct {
	Prefix string   `json:"prefix"`
	Tags   TagsList `json:"tags"`
}

// MSTeamsAccessConfig is the access configuration for Microsoft Teams profile setup.
type MSTeamsAccessConfig struct {
	MSTeamsAccessConfigInput

	BotApplication       MSTeamsBotApplication              `graphql:"botApplication"`
	PublishedApplication MSTeamsPublishedApplicationPayload `graphql:"publishedApplication"`
}

// MSTeamsBotApplication contains details about the Microsoft Teams bot application.
type MSTeamsBotApplication struct {
	DownloadURL string `graphql:"downloadURL"`
	Version     string `graphql:"version"`
}

// MSTeamsPublishedApplicationPayload contains details about the Microsoft Teams bot application publishing to the registry.
type MSTeamsPublishedApplicationPayload struct {
	Application *MSTeamsPublishedApplication `graphql:"application"`
}

// MSTeamsPublishedApplication contains details about the Microsoft Teams bot application publishing to the registry.
type MSTeamsPublishedApplication struct {
	CatalogID *string `graphql:"catalogId"`
	Version   *string `graphql:"version"`
}

// MSTeamsCustomerConfig is the customer configuration for Microsoft Teams profile setup.
type MSTeamsCustomerConfig struct {
	MSTeamsCustomerConfigInput

	RoundtripDigest string          `graphql:"roundtripDigest"`
	TerraformModule TerraformModule `graphql:"terraformModule"`
}

// MSTeamsEntityDetails has details about Microsoft Teams entities from their ID.
type MSTeamsEntityDetails struct {
	Channels []MSTeamsChannelDetail `graphql:"channels"`
	Teams    []MSTeamsTeamDetail    `graphql:"teams"`
}

// MSTeamsChannelDetail has details about a Microsoft Teams channel.
type MSTeamsChannelDetail struct {
	ID   string `graphql:"id"`
	Name string `graphql:"name"`
}

// MSTeamsTeamDetail has details about a Microsoft Teams team.
type MSTeamsTeamDetail struct {
	ID   string `graphql:"id"`
	Name string `graphql:"name"`
}

// MSTeamsChannelMapping contains mappings between IDs and target names for Microsoft Teams.
type MSTeamsChannelMapping struct {
	Name      string `graphql:"name" json:"name"`
	TeamID    UUID   `graphql:"teamId" json:"teamId"`
	ChannelID string `graphql:"channelId" json:"channelId"`
}

// MSTeamsConfigurationInput is the configuration input for the Microsoft Teams profile.
type MSTeamsConfigurationInput struct {
	AccessConfig    *MSTeamsAccessConfigInput   `json:"accessConfig"`
	ChannelMappings *[]MSTeamsChannelMapping    `json:"channelMappings"`
	CustomerConfig  *MSTeamsCustomerConfigInput `json:"customerConfig"`
}

type msteamsConfigurationInput struct {
	MSTeamsConfigurationInput

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i msteamsConfigurationInput) GetGraphQLType() string {
	return "MSTeamsConfigurationInput"
}

// JiraConfiguration is the configuration for Jira profiles.
type JiraConfiguration struct {
	URL      *string       `graphql:"url" json:"url"`
	Projects []JiraProject `graphql:"projects" json:"projects"`
	User     string        `graphql:"user" json:"user"`
	APIKey   string        `graphql:"apiKey" json:"apiKey"`
}

// JiraProject is the configuration for a Jira project.
type JiraProject struct {
	ClosedStatus string `graphql:"closedStatus" json:"closedStatus"`
	IssueType    string `graphql:"issueType" json:"issueType"`
	Name         string `graphql:"name" json:"name"`
	Project      string `graphql:"project" json:"project"`
}

type jiraConfigurationInput struct {
	JiraConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i jiraConfigurationInput) GetGraphQLType() string {
	return "JiraConfigurationInput"
}

// ResourceOwnerConfiguration is the configuration for resource owner.
type ResourceOwnerConfiguration struct {
	// "default" is present with different type in both resource and account, so it must be aliased
	Default      []string `graphql:"resourceOwnerDefault: default" json:"default"`
	OrgDomain    *string  `graphql:"orgDomain" json:"orgDomain"`
	OrgDomainTag *string  `graphql:"orgDomainTag" json:"orgDomainTag"`
	Tags         []string `graphql:"tags" json:"tags"`
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
	OrgDomain    *string         `graphql:"orgDomain" json:"orgDomain"`
	OrgDomainTag *string         `graphql:"orgDomainTag" json:"orgDomainTag"`
	Tags         []string        `graphql:"tags" json:"tags"`
}

// AccountOwners tracks the owners for an account.
type AccountOwners struct {
	Account string   `graphql:"account" json:"account"`
	Owners  []string `graphql:"owners" json:"owners"`
}

type accountOwnersConfigurationInput struct {
	AccountOwnersConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i accountOwnersConfigurationInput) GetGraphQLType() string {
	return "AccountOwnersConfigurationInput"
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

type symphonyConfigurationInput struct {
	SymphonyConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i symphonyConfigurationInput) GetGraphQLType() string {
	return "SymphonyConfigurationInput"
}

type emailConfigurationInput struct {
	EmailConfiguration

	Name  string `json:"name"`
	Scope string `json:"scope"`
}

func (i emailConfigurationInput) GetGraphQLType() string {
	return "EmailConfigurationInput"
}

type configurationProfileAPI struct {
	c *client
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
		return nil, err
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

// ReadMSTeams returns data for the Microsoft Teams configuration profile.
func (a configurationProfileAPI) ReadMSTeams(ctx context.Context) (*ConfigurationProfile, error) {
	return a.Read(ctx, ConfigurationProfileMSTeams)
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// UpsertSymphony upserts the Symphony configuration profile.
func (a configurationProfileAPI) UpsertSymphony(ctx context.Context, config SymphonyConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addSymphonyProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": symphonyConfigurationInput{
			SymphonyConfiguration: config,
			Name:                  string(ConfigurationProfileSymphony),
			Scope:                 configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// UpsertEmail creates or updates the email configuration profile.
func (a configurationProfileAPI) UpsertEmail(ctx context.Context, input EmailConfiguration) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"addEmailProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": emailConfigurationInput{
			EmailConfiguration: input,
			Name:               string(ConfigurationProfileEmail),
			Scope:              configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// UpsertMSTeams creates or updates the Microsoft Teams configuration profile.
func (a configurationProfileAPI) UpsertMSTeams(ctx context.Context, input MSTeamsConfigurationInput) (*ConfigurationProfile, error) {
	var mutation struct {
		Payload struct {
			Configuration ConfigurationProfile
		} `graphql:"upsertMSTeamsProfile(input: $input)"`
	}
	variables := map[string]any{
		"input": msteamsConfigurationInput{
			MSTeamsConfigurationInput: input,
			Name:                      string(ConfigurationProfileMSTeams),
			Scope:                     configurationScopeGlobal,
		},
	}

	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if mutation.Payload.Configuration.ID == "" {
		return nil, NotFound{"Configuration profile not found after upsert"}
	}

	return &mutation.Payload.Configuration, nil
}

// Delete removes a configuration profile.
func (a configurationProfileAPI) Delete(ctx context.Context, name ConfigurationProfileName) error {
	var mutation struct {
		ID graphql.ID `graphql:"removeProfile(scope: $scope, name: $name)"`
	}
	variables := map[string]any{
		"name":  graphql.String(string(name)),
		"scope": graphql.String(configurationScopeGlobal),
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}
	return nil
}

// DeleteJira deletes the Jira configuration profile.
func (a configurationProfileAPI) DeleteJira(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileJira)
}

// DeleteMSTeams deletes the Microsoft Teams configuration profile.
func (a configurationProfileAPI) DeleteMSTeams(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileMSTeams)
}

// DeleteSlack deletes the Slack configuration profile.
func (a configurationProfileAPI) DeleteSlack(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileSlack)
}

// DeleteServiceNow deletes the ServiceNow configuration profile.
func (a configurationProfileAPI) DeleteServiceNow(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileServiceNow)
}

// DeleteSymphony deletes the Symphony configuration profile.
func (a configurationProfileAPI) DeleteSymphony(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileSymphony)
}

// DeleteEmail deletes the email configuration profile.
func (a configurationProfileAPI) DeleteEmail(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileEmail)
}

// DeleteAccountOwners deletes the account owners configuration profile.
func (a configurationProfileAPI) DeleteAccountOwners(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileAccountOwners)
}

// DeleteResourceOwner deletes the resource owners configuration profile.
func (a configurationProfileAPI) DeleteResourceOwner(ctx context.Context) error {
	return a.Delete(ctx, ConfigurationProfileResourceOwner)
}
