// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

func TestNewStatusInfo(t *testing.T) {
	info, diags := NewStatusInfo(api.StatusInfo{
		Status: "PARTIAL_SUCCESS",
		Details: []api.StatusDetail{
			{
				Component: "us-east-1-secret",
				Status:    strPtr("SUCCESS"),
				At:        strPtr("2026-08-27T09:30:00+00:00"),
			},
			{
				Component:    "us-east-1-execution",
				Status:       strPtr("ERROR"),
				At:           strPtr("2026-08-27T09:31:00+00:00"),
				Message:      strPtr("could not read the source secret"),
				RunningSince: strPtr("2026-08-27T09:32:00+00:00"),
			},
			{
				// No completed attempt yet: a null status is not a failure.
				Component: "eu-west-1-execution",
				Message:   strPtr("the change is saved and the fan-out queued"),
			},
		},
	})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	attrs := info.Attributes()
	assert.Equal(t, types.StringValue("PARTIAL_SUCCESS"), attrs["status"])

	details, ok := attrs["details"].(types.List)
	require.True(t, ok)
	require.Len(t, details.Elements(), 3)

	first, ok := details.Elements()[0].(types.Object)
	require.True(t, ok)
	assert.Equal(t, map[string]attr.Value{
		"component":     types.StringValue("us-east-1-secret"),
		"status":        types.StringValue("SUCCESS"),
		"at":            types.StringValue("2026-08-27T09:30:00+00:00"),
		"message":       types.StringNull(),
		"running_since": types.StringNull(),
	}, first.Attributes())

	second, ok := details.Elements()[1].(types.Object)
	require.True(t, ok)
	assert.Equal(t, map[string]attr.Value{
		"component":     types.StringValue("us-east-1-execution"),
		"status":        types.StringValue("ERROR"),
		"at":            types.StringValue("2026-08-27T09:31:00+00:00"),
		"message":       types.StringValue("could not read the source secret"),
		"running_since": types.StringValue("2026-08-27T09:32:00+00:00"),
	}, second.Attributes())

	third, ok := details.Elements()[2].(types.Object)
	require.True(t, ok)
	assert.Equal(t, map[string]attr.Value{
		"component":     types.StringValue("eu-west-1-execution"),
		"status":        types.StringNull(),
		"at":            types.StringNull(),
		"message":       types.StringValue("the change is saved and the fan-out queued"),
		"running_since": types.StringNull(),
	}, third.Attributes())
}

func TestNewStatusInfo_NoDetails(t *testing.T) {
	// Nothing has reported yet, which is a known-empty list rather than an
	// unknown one, so the attribute is readable straight after an apply.
	info, diags := NewStatusInfo(api.StatusInfo{Status: "SKIPPED", Details: []api.StatusDetail{}})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	details, ok := info.Attributes()["details"].(types.List)
	require.True(t, ok)
	assert.False(t, details.IsNull())
	assert.Empty(t, details.Elements())
}
