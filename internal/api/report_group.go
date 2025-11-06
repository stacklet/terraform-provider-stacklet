// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// ReportGroup is the data for a notification report group.
type ReportGroup struct {
	ID                 string
	Name               string
	Enabled            bool
	Bindings           []string
	Source             ReportSource
	Schedule           string
	GroupBy            []string
	UseMessageSettings bool
	DeliverySettings   []struct {
		TypeName                   string                     `graphql:"__typename"`
		EmailDeliverySettings      EmailDeliverySettings      `graphql:"... on EmailSettings"`
		SlackDeliverySettings      SlackDeliverySettings      `graphql:"... on SlackSettings"`
		MSTeamsDeliverySettings    MSTeamsDeliverySettings    `graphql:"... on MSTeamsSettings"`
		ServiceNowDeliverySettings ServiceNowDeliverySettings `graphql:"... on ServiceNowSettings"`
		JiraDeliverySettings       JiraDeliverySettings       `graphql:"... on JiraSettings"`
		SymphonyDeliverySettings   SymphonyDeliverySettings   `graphql:"... on SymphonySettings"`
	}
}

// EmailDeliverySettings returns the list of email delivery settings for the
// report group.
func (r ReportGroup) EmailDeliverySettings() []EmailDeliverySettings {
	settings := make([]EmailDeliverySettings, 0)
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "EmailSettings" {
			settings = append(settings, ds.EmailDeliverySettings)
		}
	}
	return settings
}

// SlackDeliverySettings returns the list of Slack delivery settings for the
// report group.
func (r ReportGroup) SlackDeliverySettings() []SlackDeliverySettings {
	settings := make([]SlackDeliverySettings, 0)
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "SlackSettings" {
			settings = append(settings, ds.SlackDeliverySettings)
		}
	}
	return settings
}

// MSTeamsDeliverySettings returns the list of Teams delivery settings for the
// report group.
func (r ReportGroup) MSTeamsDeliverySettings() []MSTeamsDeliverySettings {
	settings := make([]MSTeamsDeliverySettings, 0)
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "MSTeamsSettings" {
			settings = append(settings, ds.MSTeamsDeliverySettings)
		}
	}
	return settings
}

// ServiceNowDeliverySettings returns the list of ServiceNow delivery settings for the
// report group.
func (r ReportGroup) ServiceNowDeliverySettings() []ServiceNowDeliverySettings {
	settings := make([]ServiceNowDeliverySettings, 0)
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "ServiceNowSettings" {
			settings = append(settings, ds.ServiceNowDeliverySettings)
		}
	}
	return settings
}

// JiraDeliverySettings returns the list of Jira delivery settings for the
// report group.
func (r ReportGroup) JiraDeliverySettings() []JiraDeliverySettings {
	settings := make([]JiraDeliverySettings, 0)
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "JiraSettings" {
			settings = append(settings, ds.JiraDeliverySettings)
		}
	}
	return settings
}

// SymphonyDeliverySettings returns the list of Symphony delivery settings for the
// report group.
func (r ReportGroup) SymphonyDeliverySettings() []SymphonyDeliverySettings {
	settings := make([]SymphonyDeliverySettings, 0)
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "SymphonySettings" {
			settings = append(settings, ds.SymphonyDeliverySettings)
		}
	}
	return settings
}

type EmailDeliverySettings struct {
	CC             []string    `json:"cc"`
	FirstMatchOnly *bool       `json:"firstMatchOnly"`
	Format         *string     `json:"format"`
	FromEmail      *string     `json:"fromEmail"`
	Priority       *string     `json:"priority"`
	Recipients     []Recipient `json:"recipients"`
	Subject        string      `json:"subject"`
	Template       string      `json:"template"`
}

type SlackDeliverySettings struct {
	FirstMatchOnly *bool       `json:"firstMatchOnly"`
	Recipients     []Recipient `json:"recipients"`
	Template       string      `json:"template"`
}

type MSTeamsDeliverySettings struct {
	FirstMatchOnly *bool       `json:"firstMatchOnly"`
	Recipients     []Recipient `json:"recipients"`
	Template       string      `json:"template"`
}

type ServiceNowDeliverySettings struct {
	FirstMatchOnly   *bool       `json:"firstMatchOnly"`
	Impact           string      `json:"impact"`
	Recipients       []Recipient `json:"recipients"`
	ShortDescription string      `json:"shortDescription"`
	Template         string      `json:"template"`
	Urgency          string      `json:"urgency"`
}

type JiraDeliverySettings struct {
	FirstMatchOnly *bool       `json:"firstMatchOnly"`
	Recipients     []Recipient `json:"recipients"`
	Template       string      `json:"template"`
	Description    string      `json:"description"`
	Project        string      `json:"project"`
	Summary        string      `json:"summary"`
}

type SymphonyDeliverySettings struct {
	FirstMatchOnly *bool       `json:"firstMatchOnly"`
	Recipients     []Recipient `json:"recipients"`
	Template       string      `json:"template"`
}

// ReportGroupsInput is the input to create or update a report group.
type ReportGroupInput struct {
	Name               string                       `json:"name"`
	Enabled            bool                         `json:"enabled"`
	Bindings           []string                     `json:"bindings"`
	Source             ReportSource                 `json:"source"`
	Schedule           string                       `json:"schedule"`
	GroupBy            []string                     `json:"groupBy"`
	UseMessageSettings bool                         `json:"useMessageSettings"`
	EmailSettings      []EmailDeliverySettings      `json:"emailSettings"`
	SlackSettings      []SlackDeliverySettings      `json:"slackSettings"`
	MSTeamsSettings    []MSTeamsDeliverySettings    `json:"msteamsSettings"`
	ServiceNowSettings []ServiceNowDeliverySettings `json:"serviceNowSettings"`
	JiraSettings       []JiraDeliverySettings       `json:"jiraSettings"`
	SymphonySettings   []SymphonyDeliverySettings   `json:"symphonySettings"`
}

type Recipient struct {
	AccountOwner  *bool   `graphql:"account_owner" json:"account_owner"`
	EventOwner    *bool   `graphql:"event_owner" json:"event_owner"`
	ResourceOwner *bool   `graphql:"resource_owner" json:"resource_owner"`
	Tag           *string `json:"tag"`
	Value         *string `json:"value"`
}

type upsertReportGroupsInput struct {
	ReportGroups []ReportGroupInput `json:"reportGroups"`
}

func (i upsertReportGroupsInput) GetGraphQLType() string {
	return "UpsertReportGroupsInput"
}

type reportGroupAPI struct {
	c *graphql.Client
}

// Read returns data for a report group.
func (a reportGroupAPI) Read(ctx context.Context, name string) (*ReportGroup, error) {
	var query struct {
		ReportGroup ReportGroup `graphql:"reportGroup(name: $name)"`
	}

	if err := a.c.Query(ctx, &query, map[string]any{"name": graphql.String(name)}); err != nil {
		return nil, NewAPIError(err)
	}

	if query.ReportGroup.ID == "" {
		return nil, NotFound{"Report group not found"}
	}

	return &query.ReportGroup, nil
}

// Upsert creates or updates a report group.
func (a reportGroupAPI) Upsert(ctx context.Context, input ReportGroupInput) (*ReportGroup, error) {
	var mutation struct {
		Payload struct {
			ReportGroups []ReportGroup
		} `graphql:"upsertReportGroups(input: $input)"`
	}
	variables := map[string]any{
		"input": upsertReportGroupsInput{
			ReportGroups: []ReportGroupInput{input},
		},
	}
	if err := a.c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, NewAPIError(err)
	}

	if len(mutation.Payload.ReportGroups) == 0 {
		return nil, NotFound{"Report group not found after upsert"}
	}
	return &mutation.Payload.ReportGroups[0], nil
}

// Delete removes a report group.
func (a reportGroupAPI) Delete(ctx context.Context, name string) error {
	var mutation struct {
		IDs []string `graphql:"removeReportGroups(names: $names)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"names": []graphql.String{graphql.String(name)}}); err != nil {
		return NewAPIError(err)
	}
	return nil
}
