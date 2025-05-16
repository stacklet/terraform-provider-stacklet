package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// AddDiagError adds an error to the diagnostics.
func AddDiagError(diag *diag.Diagnostics, err error) {
	switch e := err.(type) {
	case api.APIError:
		diag.AddError(e.Summary, e.Detail)
	case ImportIDError:
		diag.AddError(e.Summary(), e.Error())
	default:
		diag.AddError("Error", e.Error())
	}
}
