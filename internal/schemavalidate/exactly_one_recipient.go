// Copyright (c) 2025 - Stacklet, Inc.

package schemavalidate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExactlyOneRecipient returns a validator that ensures that exactly one field
// is set in for recipients: either one of the owners flags are true, or
// exactly one of tag and value is set.
func ExactlyOneRecipient() validator.Object {
	return exactlyOneRecipient{}
}

type exactlyOneRecipient struct{}

func (v exactlyOneRecipient) Description(ctx context.Context) string {
	return "Ensures exactly one recipient field is set: either one of the owners flags are true, or exactly one of tag and value is set"
}

func (v exactlyOneRecipient) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v exactlyOneRecipient) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	setCount := 0
	for _, attr := range req.ConfigValue.Attributes() {
		if attr.IsNull() || attr.IsUnknown() {
			continue
		}
		switch a := attr.(type) {
		case types.Bool:
			if a.ValueBool() {
				setCount++
			}
		case types.String:
			if a.ValueString() != "" {
				setCount++
			}
		}
	}

	if setCount != 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid recipient configuration",
			"Exactly one recipient field must be set.",
		)
		return
	}
}
