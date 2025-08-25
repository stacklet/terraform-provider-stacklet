// Copyright (c) 2025 - Stacklet, Inc.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DefaultStringListEmpty is a defaults.List default for an empty list of strings.
func DefaultStringListEmpty() defaults.List {
	return listdefault.StaticValue(basetypes.NewListValueMust(types.StringType, []attr.Value{}))
}
