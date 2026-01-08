// Copyright (c) 2025 - Stacklet, Inc.

package planmodifiers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestTrimWhiteSpace(t *testing.T) {
	tests := []struct {
		name      string
		planValue types.String
		expected  types.String
	}{
		{
			name:      "single trailing newline",
			planValue: types.StringValue("content\n"),
			expected:  types.StringValue("content"),
		},
		{
			name:      "multiple trailing newlines",
			planValue: types.StringValue("content\n\n\n"),
			expected:  types.StringValue("content"),
		},
		{
			name:      "trailing carriage return and newline",
			planValue: types.StringValue("content\r\n"),
			expected:  types.StringValue("content"),
		},
		{
			name:      "mixed trailing whitespace",
			planValue: types.StringValue("content\n\r\n\r"),
			expected:  types.StringValue("content"),
		},
		{
			name:      "no trailing newlines",
			planValue: types.StringValue("content"),
			expected:  types.StringValue("content"),
		},
		{
			name:      "newlines in middle preserved",
			planValue: types.StringValue("line1\nline2\nline3\n"),
			expected:  types.StringValue("line1\nline2\nline3"),
		},
		{
			name:      "empty string",
			planValue: types.StringValue(""),
			expected:  types.StringValue(""),
		},
		{
			name:      "only newlines",
			planValue: types.StringValue("\n\n\n"),
			expected:  types.StringValue(""),
		},
		{
			name:      "null value unchanged",
			planValue: types.StringNull(),
			expected:  types.StringNull(),
		},
		{
			name:      "unknown value unchanged",
			planValue: types.StringUnknown(),
			expected:  types.StringUnknown(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			req := planmodifier.StringRequest{
				PlanValue: tt.planValue,
			}
			resp := &planmodifier.StringResponse{
				PlanValue: tt.planValue,
			}

			modifier := TrimWhiteSpace()
			modifier.PlanModifyString(ctx, req, resp)

			assert.Equal(t, tt.expected, resp.PlanValue)
		})
	}
}

func TestTrimWhiteSpace_Metadata(t *testing.T) {
	modifier := TrimWhiteSpace()
	ctx := t.Context()

	description := modifier.Description(ctx)
	assert.NotEmpty(t, description)
	assert.Contains(t, description, "trailing whitespace")

	markdownDescription := modifier.MarkdownDescription(ctx)
	assert.NotEmpty(t, markdownDescription)
	assert.Contains(t, markdownDescription, "trailing whitespace")
}

func TestRequiresReplaceIfNullStringChange(t *testing.T) {
	tests := []struct {
		name            string
		stateValue      types.String
		planValue       types.String
		stateRawIsNull  bool
		requiresReplace bool
	}{
		{
			name:            "no change - both non-null",
			stateValue:      types.StringValue("value1"),
			planValue:       types.StringValue("value2"),
			stateRawIsNull:  false,
			requiresReplace: false,
		},
		{
			name:            "no change - both null",
			stateValue:      types.StringNull(),
			planValue:       types.StringNull(),
			stateRawIsNull:  false,
			requiresReplace: false,
		},
		{
			name:            "null to non-null",
			stateValue:      types.StringNull(),
			planValue:       types.StringValue("value"),
			stateRawIsNull:  false,
			requiresReplace: true,
		},
		{
			name:            "non-null to null",
			stateValue:      types.StringValue("value"),
			planValue:       types.StringNull(),
			stateRawIsNull:  false,
			requiresReplace: true,
		},
		{
			name:            "create operation - skip check",
			stateValue:      types.StringNull(),
			planValue:       types.StringValue("value"),
			stateRawIsNull:  true,
			requiresReplace: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			req := planmodifier.StringRequest{
				StateValue: tt.stateValue,
				PlanValue:  tt.planValue,
				State: tfsdk.State{
					Raw: createRawValue(tt.stateRawIsNull),
				},
			}
			resp := &planmodifier.StringResponse{
				RequiresReplace: false,
			}

			modifier := RequiresReplaceIfNullStringChange()
			modifier.PlanModifyString(ctx, req, resp)

			assert.Equal(t, tt.requiresReplace, resp.RequiresReplace)
		})
	}
}

func TestRequiresReplaceIfNullStringChange_Metadata(t *testing.T) {
	modifier := RequiresReplaceIfNullStringChange()
	ctx := t.Context()

	description := modifier.Description(ctx)
	assert.NotEmpty(t, description)
	assert.Contains(t, description, "null")

	markdownDescription := modifier.MarkdownDescription(ctx)
	assert.NotEmpty(t, markdownDescription)
	assert.Contains(t, markdownDescription, "null")
}

// createRawValue creates a mock raw value for testing.
func createRawValue(isNull bool) tftypes.Value {
	if isNull {
		return tftypes.NewValue(tftypes.Object{}, nil)
	}
	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"test": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"test": tftypes.NewValue(tftypes.String, "value"),
	})
}
