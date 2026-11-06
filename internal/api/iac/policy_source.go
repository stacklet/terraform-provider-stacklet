// Copyright Stacklet, Inc. 2025, 2026

package iac

import (
	"context"
	"encoding/json"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// PolicySourceAuth is the authentication for a policy source.
//
// Either Token, or SSHKey (with optional SSHPassphrase) is set. When returned
// by the API, secret values are masked.
type PolicySourceAuth struct {
	Username      string  `json:"username"`
	Token         *string `json:"token,omitempty"`
	SSHKey        *string `json:"ssh_key,omitempty"`
	SSHPassphrase *string `json:"ssh_passphrase,omitempty"`
}

func (a *PolicySourceAuth) secrets() []*string {
	if a == nil {
		return nil
	}
	return []*string{a.Token, a.SSHKey, a.SSHPassphrase}
}

// PolicySource is a repository of policies.
type PolicySource struct {
	Name       string            `json:"name"`
	URL        string            `json:"url"`
	Branch     *string           `json:"branch"`
	Paths      []string          `json:"paths"`
	Extensions []string          `json:"extensions"`
	AutoScan   bool              `json:"auto_scan"`
	HookPath   string            `json:"hook_path"`
	Auth       *PolicySourceAuth `json:"auth"`

	// WebhookURL is the full URL for the webhook, built from HookPath.
	WebhookURL string `json:"-"`
}

// PolicySourceCreateInput is the input for creating a policy source.
type PolicySourceCreateInput struct {
	Name       string            `json:"name"`
	URL        string            `json:"url"`
	Branch     *string           `json:"branch,omitempty"`
	Paths      []string          `json:"paths,omitempty"`
	Extensions []string          `json:"extensions,omitempty"`
	AutoScan   bool              `json:"auto_scan"`
	Auth       *PolicySourceAuth `json:"auth,omitempty"`
}

// PolicySourceAuthUpdate is an update to policy source authentication.
//
// The API replaces authentication as a whole. If Set is false, authentication
// is left unchanged, otherwise it's replaced with Auth, or removed if that's nil.
type PolicySourceAuthUpdate struct {
	Set  bool
	Auth *PolicySourceAuth
}

func (u PolicySourceAuthUpdate) IsZero() bool {
	return !u.Set
}

func (u PolicySourceAuthUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Auth)
}

// PolicySourceUpdateInput is the input for updating a policy source.
type PolicySourceUpdateInput struct {
	URL        string                 `json:"url"`
	Branch     *string                `json:"branch"`
	Paths      []string               `json:"paths"`
	Extensions []string               `json:"extensions,omitempty"`
	AutoScan   bool                   `json:"auto_scan"`
	Auth       PolicySourceAuthUpdate `json:"auth,omitzero"`
}

type policySourceAPI struct {
	c *client
}

func (a policySourceAPI) Read(ctx context.Context, name string) (*PolicySource, error) {
	var source PolicySource
	if err := a.c.Get(ctx, resourcePath("policy-sources", name), nil, &source); err != nil {
		return nil, err
	}
	return a.withWebhookURL(&source)
}

func (a policySourceAPI) Create(ctx context.Context, i PolicySourceCreateInput) (*PolicySource, error) {
	ctx = api.MaskSecrets(ctx, i.Auth.secrets()...)
	var source PolicySource
	if err := a.c.Post(ctx, "policy-sources", i, &source); err != nil {
		return nil, err
	}
	return a.withWebhookURL(&source)
}

func (a policySourceAPI) Update(ctx context.Context, name string, i PolicySourceUpdateInput) (*PolicySource, error) {
	ctx = api.MaskSecrets(ctx, i.Auth.Auth.secrets()...)
	if err := a.c.Patch(ctx, resourcePath("policy-sources", name), i, nil); err != nil {
		return nil, err
	}
	// The update response doesn't include authentication details, so read
	// the policy source back.
	return a.Read(ctx, name)
}

func (a policySourceAPI) Delete(ctx context.Context, name string) error {
	return a.c.Delete(ctx, resourcePath("policy-sources", name), nil)
}

func (a policySourceAPI) withWebhookURL(source *PolicySource) (*PolicySource, error) {
	webhookURL, err := a.c.URL(source.HookPath)
	if err != nil {
		return nil, err
	}
	source.WebhookURL = webhookURL
	return source, nil
}
