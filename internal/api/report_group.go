// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

// ReportGroup is the data for a notification report group.
type ReportGroup struct {
	ID                 graphql.ID   `graphql:"id"`
	Name               string       `graphql:"name"`
	Enabled            bool         `graphql:"enabled"`
	Bindings           []string     `graphql:"bindings"`
	Source             ReportSource `graphql:"source"`
	Schedule           string       `graphql:"schedule"`
	GroupBy            []string     `graphql:"groupBy"`
	UseMessageSettings bool         `graphql:"useMessageSettings"`
	DeliverySettings   []struct {
		TypeName                   string                     `graphql:"__typename"`
		EmailDeliverySettings      EmailDeliverySettings      `graphql:"... on EmailSettings"`
		SlackDeliverySettings      SlackDeliverySettings      `graphql:"... on SlackSettings"`
		MSTeamsDeliverySettings    MSTeamsDeliverySettings    `graphql:"... on MSTeamsSettings"`
		ServiceNowDeliverySettings ServiceNowDeliverySettings `graphql:"... on ServiceNowSettings"`
		JiraDeliverySettings       JiraDeliverySettings       `graphql:"... on JiraSettings"`
		SymphonyDeliverySettings   SymphonyDeliverySettings   `graphql:"... on SymphonySettings"`
	} `graphql:"deliverySettings"`
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
	CC             []string    `graphql:"cc" json:"cc"`
	FirstMatchOnly *bool       `graphql:"firstMatchOnly" json:"firstMatchOnly"`
	Format         *string     `graphql:"format" json:"format"`
	FromEmail      *string     `graphql:"fromEmail" json:"fromEmail"`
	Priority       *string     `graphql:"priority" json:"priority"`
	Recipients     []Recipient `graphql:"recipients" json:"recipients"`
	Subject        string      `graphql:"subject" json:"subject"`
	Template       string      `graphql:"template" json:"template"`
}

type SlackDeliverySettings struct {
	FirstMatchOnly *bool       `graphql:"firstMatchOnly" json:"firstMatchOnly"`
	Recipients     []Recipient `graphql:"recipients" json:"recipients"`
	Template       string      `graphql:"template" json:"template"`
}

type MSTeamsDeliverySettings struct {
	FirstMatchOnly *bool       `graphql:"firstMatchOnly" json:"firstMatchOnly"`
	Recipients     []Recipient `graphql:"recipients" json:"recipients"`
	Template       string      `graphql:"template" json:"template"`
}

type ServiceNowDeliverySettings struct {
	FirstMatchOnly   *bool       `graphql:"firstMatchOnly" json:"firstMatchOnly"`
	Impact           string      `graphql:"impact" json:"impact"`
	Recipients       []Recipient `graphql:"recipients" json:"recipients"`
	ShortDescription string      `graphql:"shortDescription" json:"shortDescription"`
	Template         string      `graphql:"template" json:"template"`
	Urgency          string      `graphql:"urgency" json:"urgency"`
}

type JiraDeliverySettings struct {
	FirstMatchOnly *bool       `graphql:"firstMatchOnly" json:"firstMatchOnly"`
	Recipients     []Recipient `graphql:"recipients" json:"recipients"`
	Template       string      `graphql:"template" json:"template"`
	Description    string      `graphql:"description" json:"description"`
	Project        string      `graphql:"project" json:"project"`
	Summary        string      `graphql:"summary" json:"summary"`
}

type SymphonyDeliverySettings struct {
	FirstMatchOnly *bool       `graphql:"firstMatchOnly" json:"firstMatchOnly"`
	Recipients     []Recipient `graphql:"recipients" json:"recipients"`
	Template       string      `graphql:"template" json:"template"`
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
	Tag           *string `graphql:"tag" json:"tag"`
	Value         *string `graphql:"value" json:"value"`
}

type upsertReportGroupsInput struct {
	ReportGroups []ReportGroupInput `json:"reportGroups"`
}

func (i upsertReportGroupsInput) GetGraphQLType() string {
	return "UpsertReportGroupsInput"
}

type reportGroupAPI struct {
	c *client
}

// Read returns data for a report group.
func (a reportGroupAPI) Read(ctx context.Context, name string) (*ReportGroup, error) {
	var query struct {
		ReportGroup ReportGroup `graphql:"reportGroup(name: $name)"`
	}

	if err := a.c.Query(ctx, &query, map[string]any{"name": graphql.String(name)}); err != nil {
		return nil, err
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
		return nil, err
	}

	if len(mutation.Payload.ReportGroups) == 0 {
		return nil, NotFound{"Report group not found after upsert"}
	}
	return &mutation.Payload.ReportGroups[0], nil
}

// Delete removes a report group.
func (a reportGroupAPI) Delete(ctx context.Context, name string) error {
	var mutation struct {
		IDs []graphql.ID `graphql:"removeReportGroups(names: $names)"`
	}
	if err := a.c.Mutate(ctx, &mutation, map[string]any{"names": []graphql.String{graphql.String(name)}}); err != nil {
		return err
	}
	return nil
}
