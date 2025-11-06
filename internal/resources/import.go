// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

type importIDError struct {
	parts []string
}

// Error returns the error message.
func (e importIDError) Error() string {
	format := strings.Join(e.parts, ":")
	return fmt.Sprintf("Import ID must be in the format: %s", format)
}

// Summary returns a summary message for the error.
func (e importIDError) Summary() string {
	return "Invalid import ID"
}

// splitImportID splits an import ID into expected components.
func splitImportID(id string, parts []string) ([]string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != len(parts) {
		return nil, importIDError{parts: parts}
	}
	return idParts, nil
}

// importState imports the state identifier from the request ID. If multiple attributes are specified, it's assumed the ID is built by concatenating them vith `:`.
func importState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, attrs []string) {
	var values []string

	if len(attrs) == 1 {
		values = []string{req.ID}
	} else {
		var err error
		values, err = splitImportID(req.ID, attrs)
		if err != nil {
			errors.AddDiagError(&resp.Diagnostics, err)
			return
		}
	}

	for i, attr := range attrs {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attr), values[i])...)
	}
}
