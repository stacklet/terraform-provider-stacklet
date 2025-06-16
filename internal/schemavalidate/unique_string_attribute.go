// Copyright (c) 2025 - Stacklet, Inc.

package schemavalidate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UniqueStringAttribute returns a validator that ensures that each
// element in a list has different values for the specified string attribute.
func UniqueStringAttribute(name string) validator.List {
	return uniqueStringAttribute{name: name}
}

type uniqueStringAttribute struct {
	name string
}

func (v uniqueStringAttribute) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures all entries have unique values for string attribute '%s'", v.name)
}

func (v uniqueStringAttribute) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uniqueStringAttribute) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	seen := make(map[string]bool)

	for i, element := range elements {
		if element.IsNull() || element.IsUnknown() {
			continue
		}

		obj, ok := element.(types.Object)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i),
				"Validation Error",
				"Element is not a types.Object",
			)
			return
		}
		attrs := obj.Attributes()
		attr, exists := attrs[v.name]
		if !exists {
			continue
		}

		a, ok := attr.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i).AtName(v.name),
				"Validation Error",
				"Element attribute is not of types.String",
			)
			return
		}

		value := a.ValueString()
		if seen[value] {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i).AtName(v.name),
				"Duplicate Value",
				fmt.Sprintf("Value '%s' for attribute '%s' must be unique across all entries", value, v.name),
			)
			return
		}
		seen[value] = true
	}
}
