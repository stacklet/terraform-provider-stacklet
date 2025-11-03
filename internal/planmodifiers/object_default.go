// Copyright (c) 2025 - Stacklet, Inc.

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DefaultObject returns a plan modifier that sets a default object value when
// the planned value is null. This is useful when the API applies defaults that
// need to be reflected in the plan to avoid consistency errors.
func DefaultObject(defaultValue basetypes.ObjectValue) planmodifier.Object {
	return defaultObjectModifier{
		defaultValue: defaultValue,
	}
}

type defaultObjectModifier struct {
	defaultValue basetypes.ObjectValue
}

func (m defaultObjectModifier) Description(ctx context.Context) string {
	return "Sets a default value for the object if the planned value is null."
}

func (m defaultObjectModifier) MarkdownDescription(ctx context.Context) string {
	return "Sets a default value for the object if the planned value is null."
}

func (m defaultObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If the configuration is null and the plan is null, set the plan to the default value
	if req.ConfigValue.IsNull() && req.PlanValue.IsNull() {
		resp.PlanValue = m.defaultValue
	}
}
