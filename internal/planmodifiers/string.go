// Copyright (c) 2025 - Stacklet, Inc.

package planmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RequiresReplaceIfNullStringChange returns a planmodifier.String that requires replacement
// when a string attribute transitions between null and non-null (in either direction).
func RequiresReplaceIfNullStringChange() planmodifier.String {
	return requiresReplaceIfNullStringChange{}
}

type requiresReplaceIfNullStringChange struct{}

func (m requiresReplaceIfNullStringChange) Description(ctx context.Context) string {
	return "Requires replacement when transitioning between null and non-null."
}

func (m requiresReplaceIfNullStringChange) MarkdownDescription(ctx context.Context) string {
	return "Requires replacement when transitioning between null and non-null."
}

func (m requiresReplaceIfNullStringChange) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Skip on create
	if req.State.Raw.IsNull() {
		return
	}

	if req.StateValue.IsNull() != req.PlanValue.IsNull() {
		resp.RequiresReplace = true
	}
}

// TrimWhiteSpace returns a planmodifier.String that removes leading and trailing unicode whitespace characters
// from the planned value. This is useful when an API normalizes content by removing
// trailing whitespace, preventing spurious diffs.
func TrimWhiteSpace() planmodifier.String {
	return trimWhiteSpace{}
}

type trimWhiteSpace struct{}

func (m trimWhiteSpace) Description(ctx context.Context) string {
	return "Trims leading and trailing whitespace from the planned value to match API normalization."
}

func (m trimWhiteSpace) MarkdownDescription(ctx context.Context) string {
	return "Trims leading and trailing whitespace from the planned value to match API normalization."
}

func (m trimWhiteSpace) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the plan value is null or unknown, don't modify it
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// Trim leading and trailing whitespace from the plan value
	trimmed := strings.TrimSpace(req.PlanValue.ValueString())
	resp.PlanValue = types.StringValue(trimmed)
}
