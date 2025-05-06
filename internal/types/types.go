package types

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullableString return the proper type for a nullable string
func NullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// JSONMap converts a JSON-encoded string to a map
func JSONMap(ctx context.Context, s string) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	if s == "" {
		return types.MapNull(types.StringType), diags
	}

	var rawMap map[string]string
	if err := json.Unmarshal([]byte(s), &rawMap); err != nil {
		diags.AddError(
			"Invalid JSON Format",
			fmt.Sprintf("Unable to parse JSON string: %s", err.Error()),
		)
		return types.MapNull(types.StringType), diags
	}

	tfMap := make(map[string]types.String)
	for key, value := range rawMap {
		tfMap[key] = types.StringValue(value)
	}

	result, mapDiags := types.MapValueFrom(ctx, types.StringType, tfMap)
	diags.Append(mapDiags...)

	return result, diags
}
