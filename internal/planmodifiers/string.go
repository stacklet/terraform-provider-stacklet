// Copyright (c) 2025 - Stacklet, Inc.

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
