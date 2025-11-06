// Copyright (c) 2025 - Stacklet, Inc.

package errors

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// DiagError represents an error that gets reported in a Diagnostic.
type DiagError interface {
	error

	// Summary returns the error summary
	Summary() string
}

// diagError wraps an error as a DiagError.
type diagError struct {
	err error
}

func (e diagError) Error() string {
	return e.err.Error()
}

func (e diagError) Summary() string {
	return "Error"
}

// AsDiagError ensures an error matches the DiagError interface.
func AsDiagError(err error) DiagError {
	e, ok := err.(DiagError)
	if ok {
		return e
	}
	return diagError{err: e}
}

// AddDiagError adds an error to the diagnostics.
func AddDiagError(diag *diag.Diagnostics, err error) {
	e := AsDiagError(err)
	diag.AddError(e.Summary(), e.Error())
}

// AddAttributeDiags adds attribute-specific Diagnostics to a set of Diagnostics.
func AddAttributeDiags(dest *diag.Diagnostics, src diag.Diagnostics, attr string) {
	for _, d := range src {
		dest.Append(diag.WithPath(path.Root(attr), d))
	}
}
