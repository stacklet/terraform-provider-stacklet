// Copyright (c) 2025 - Stacklet, Inc.

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfFieldsChanged returns a planmodifier.Object that causes the
// resource to be replaced value for specified fields changes.
func RequiresReplaceIfFieldsChanged(names ...string) planmodifier.Object {
	return requiresReplaceIfFieldsChanged{fieldNames: names}
}

type requiresReplaceIfFieldsChanged struct {
	fieldNames []string
}

func (m requiresReplaceIfFieldsChanged) Description(ctx context.Context) string {
	return "Requires replace if value for fields is changed or the object is removed."
}

func (m requiresReplaceIfFieldsChanged) MarkdownDescription(ctx context.Context) string {
	return "Requires replace if value for fields is changed or the object is removed."
}

func (m requiresReplaceIfFieldsChanged) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Always allow creation
	if req.State.Raw.IsNull() {
		return
	}

	stateAttrs := req.StateValue.Attributes()
	planAttrs := req.PlanValue.Attributes()
	for _, name := range m.fieldNames {
		planValue, planFound := planAttrs[name]
		stateValue, stateFound := stateAttrs[name]
		// if the object was or will be null, attributes might not be found, thus their value would be nil
		if (planFound != stateFound) || (planFound && stateFound && !planValue.Equal(stateValue)) {
			resp.RequiresReplace = true
			return
		}
	}
}

// RequiresReplaceIfUnset returns a planmodifier.Object that requires replacement if
// a field is set to null after being set to a non-null value.
func RequiresReplaceIfUnset() planmodifier.Object {
	return requiresReplaceIfUnset{}
}

type requiresReplaceIfUnset struct{}

func (m requiresReplaceIfUnset) Description(ctx context.Context) string {
	return "Requires replace if the field is removed or set to null after being set to a non-null value."
}

func (m requiresReplaceIfUnset) MarkdownDescription(ctx context.Context) string {
	return "Requires replace if the field is removed or set to null after being set to a non-null value."
}

func (m requiresReplaceIfUnset) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Allow resource creation
	if req.State.Raw.IsNull() {
		return
	}

	// If state has a non-null value but config is trying to set it to null, require replacement
	// We check the config value (not plan value) because plan modifiers may have filled in defaults
	if !req.StateValue.IsNull() && req.ConfigValue.IsNull() {
		resp.RequiresReplace = true
	}
}
