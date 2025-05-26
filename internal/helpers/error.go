package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

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
	switch e := err.(type) {
	case api.APIError:
		diag.AddError(e.Summary, e.Detail)
	case api.NotFound:
		diag.AddError(e.Summary(), e.Error())
	case providerdata.ProviderDataError:
		diag.AddError(e.Summary(), e.Error())
	case ImportIDError:
		diag.AddError(e.Summary(), e.Error())
	default:
		diag.AddError("Error", e.Error())
	}
}
