// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// EmtpyListDefault returns an empty default for a resource field.
func EmptyListDefault(attrType attr.Type) defaults.List {
	return listdefault.StaticValue(basetypes.NewListValueMust(attrType, []attr.Value{}))
}
