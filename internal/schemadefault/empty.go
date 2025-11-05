// Copyright (c) 2025 - Stacklet, Inc.

package schemadefault

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EmtpyListDefault returns an empty default for a resource field.
func EmptyListDefault(attrType attr.Type) defaults.List {
	return listdefault.StaticValue(types.ListValueMust(attrType, []attr.Value{}))
}

// EmptyMapDefault returns an empty map default for a resource field.
func EmptyMapDefault(elementType attr.Type) defaults.Map {
	return mapdefault.StaticValue(types.MapValueMust(elementType, map[string]attr.Value{}))
}
