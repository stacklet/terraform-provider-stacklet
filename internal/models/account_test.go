// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

func strPtr(s string) *string { return &s }

func TestAccountResourceUpdate_VariablesNullWhenAPIReturnsEmptyObject(t *testing.T) {
	// Simulate a user who has not set variables in their config (null in plan/state).
	// The API returns "{}" as its default empty representation.
	// After Update, variables must remain null to avoid "inconsistent result after apply".
	m := &AccountResource{}
	m.Variables = types.StringNull()

	account := &api.Account{
		ID:        "id-1",
		Key:       "123456789012",
		Name:      "test",
		Provider:  api.CloudProvider("AWS"),
		Variables: strPtr("{}"),
	}

	diags := m.Update(account)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if !m.Variables.IsNull() {
		t.Errorf("expected Variables to be null when API returns {}, got %q", m.Variables.ValueString())
	}
}

func TestAccountResourceUpdate_VariablesPreservedWhenSet(t *testing.T) {
	// When the user has explicitly set variables, the value from the API should be used.
	m := &AccountResource{}
	m.Variables = types.StringValue(`{"env":"prod"}`)

	account := &api.Account{
		ID:        "id-1",
		Key:       "123456789012",
		Name:      "test",
		Provider:  api.CloudProvider("AWS"),
		Variables: strPtr(`{"env":"prod"}`),
	}

	diags := m.Update(account)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if m.Variables.ValueString() != `{"env":"prod"}` {
		t.Errorf("expected Variables to be %q, got %q", `{"env":"prod"}`, m.Variables.ValueString())
	}
}

func TestAccountResourceUpdate_VariablesNullWhenAPIReturnsNil(t *testing.T) {
	m := &AccountResource{}
	m.Variables = types.StringNull()

	account := &api.Account{
		ID:       "id-1",
		Key:      "123456789012",
		Name:     "test",
		Provider: api.CloudProvider("AWS"),
		// Variables is nil
	}

	diags := m.Update(account)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if !m.Variables.IsNull() {
		t.Errorf("expected Variables to be null when API returns nil, got %q", m.Variables.ValueString())
	}
}
