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
		TeamsDeliverySettings      TeamsDeliverySettings      `graphql:"... on TeamsSettings"`
		ServiceNowDeliverySettings ServiceNowDeliverySettings `graphql:"... on ServiceNowSettings"`
		JiraDeliverySettings       JiraDeliverySettings       `graphql:"... on JiraSettings"`
		SymphonyDeliverySettings   SymphonyDeliverySettings   `graphql:"... on SymphonySettings"`
	}
}

// EmailDeliverySettings returns the list of email delivery settings for the
// report group.
func (r ReportGroup) EmailDeliverySettings() []EmailDeliverySettings {
	settings := []EmailDeliverySettings{}
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
	settings := []SlackDeliverySettings{}
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "SlackSettings" {
			settings = append(settings, ds.SlackDeliverySettings)
		}
	}
	return settings
}

// TeamsDeliverySettings returns the list of Teams delivery settings for the
// report group.
func (r ReportGroup) TeamsDeliverySettings() []TeamsDeliverySettings {
	settings := []TeamsDeliverySettings{}
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "TeamsSettings" {
			settings = append(settings, ds.TeamsDeliverySettings)
		}
	}
	return settings
}

// ServiceNowDeliverySettings returns the list of ServiceNow delivery settings for the
// report group.
func (r ReportGroup) ServiceNowDeliverySettings() []ServiceNowDeliverySettings {
	settings := []ServiceNowDeliverySettings{}
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
	settings := []JiraDeliverySettings{}
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
	settings := []SymphonyDeliverySettings{}
	for _, ds := range r.DeliverySettings {
		if ds.TypeName == "SymphonySettings" {
			settings = append(settings, ds.SymphonyDeliverySettings)
		}
	}
	return settings
}

type EmailDeliverySettings struct {
	CC             []string
	FirstMatchOnly *bool
	Format         *string
	FromEmail      *string
	Priority       *string
	Recipients     []Recipient
	Subject        string
	Template       string
}

type SlackDeliverySettings struct {
	FirstMatchOnly *bool
	Recipients     []Recipient
	Template       string
}

type TeamsDeliverySettings struct {
	FirstMatchOnly *bool
	Recipients     []Recipient
	Template       string
}

type ServiceNowDeliverySettings struct {
	FirstMatchOnly   *bool
	Impact           string
	Recipients       []Recipient
	ShortDescription string
	Template         string
	Urgency          string
}

type JiraDeliverySettings struct {
	FirstMatchOnly *bool
	Recipients     []Recipient
	Template       string
	Description    string
	Project        string
	Summary        string
}

type SymphonyDeliverySettings struct {
	FirstMatchOnly *bool
	Recipients     []Recipient
	Template       string
}

type Recipient struct {
	AccountOwner  *bool `graphql:"account_owner"`
	EventOwner    *bool `graphql:"event_owner"`
	ResourceOwner *bool `graphql:"resource_owner"`
	Tag           *string
	Value         *string
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
