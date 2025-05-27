// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"fmt"
	"strings"
)

// ImportIDError represent an error importing an ID.
type ImportIDError struct {
	parts []string
}

// Error returns the error message.
func (e ImportIDError) Error() string {
	format := strings.Join(e.parts, ":")
	return fmt.Sprintf("Import ID must be in the format: %s", format)
}

// Summary returns a summary message for the error.
func (e ImportIDError) Summary() string {
	return "Invalid import ID"
}

// splitImportID splits an import ID into expected components.
func splitImportID(id string, parts []string) ([]string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != len(parts) {
		return nil, ImportIDError{parts: parts}
	}
	return idParts, nil
}
