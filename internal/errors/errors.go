package errors

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// DiagError represents an error that gets reported in a Diagnostic.
type DiagError interface {
	error

	// Summary returns the error summary
	Summary() string
}

// AddDiagError adds an error to the diagnostics.
func AddDiagError(diag *diag.Diagnostics, err error) {
	if e, ok := err.(DiagError); ok {
		diag.AddError(e.Summary(), e.Error())
	} else {
		diag.AddError("Error", e.Error())
	}
}
