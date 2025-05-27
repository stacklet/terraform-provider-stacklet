package errors

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// DiagError represents an error that gets reported in a Diagnostic.
type DiagError interface {
	error

	// Summary returns the error summary
	Summary() string
}

// HandleAPIError handles errors returned from the API.
func HandleAPIError(ctx context.Context, state *tfsdk.State, diag *diag.Diagnostics, err error) {
	if _, ok := err.(api.NotFound); ok {
		state.RemoveResource(ctx)
	} else {
		AddDiagError(diag, err)
	}
}

// AddDiagError adds an error to the diagnostics.
func AddDiagError(diag *diag.Diagnostics, err error) {
	if e, ok := err.(DiagError); ok {
		diag.AddError(e.Summary(), e.Error())
	} else {
		diag.AddError("Error", e.Error())
	}
}
