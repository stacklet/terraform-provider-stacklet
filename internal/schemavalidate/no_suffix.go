// Copyright Stacklet, Inc. 2025, 2026

package schemavalidate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NoSuffix returns a validator that ensures a string doesn't end with the
// specified suffix.
func NoSuffix(suffix string) validator.String {
	return noSuffix{suffix: suffix}
}

type noSuffix struct {
	suffix string
}

func (v noSuffix) Description(ctx context.Context) string {
	return fmt.Sprintf("value must not end with %q", v.suffix)
}

func (v noSuffix) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v noSuffix) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if strings.HasSuffix(req.ConfigValue.ValueString(), v.suffix) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Value",
			fmt.Sprintf("Attribute %s %s, got: %s", req.Path, v.Description(ctx), req.ConfigValue.ValueString()),
		)
	}
}
