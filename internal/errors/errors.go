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

// AddDiagAttributeError adds an error to the diagnostics for a specific attribute.
func AddDiagAttributeError(diag *diag.Diagnostics, attr string, err error) {
	e := AsDiagError(err)
	diag.AddAttributeError(path.Root(attr), e.Summary(), e.Error())
}
