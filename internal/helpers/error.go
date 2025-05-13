package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// AddErrors adds an error to the diagnostics.
func AddDiagError(diag diag.Diagnostics, err error) {
	switch e := err.(type) {
	case api.APIError:
		diag.AddError(e.Summary, e.Detail)
	default:
		diag.AddError("Error", e.Error())
	}
}
