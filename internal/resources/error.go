package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

// handleAPIError handles errors returned from the API.
func handleAPIError(ctx context.Context, state *tfsdk.State, diag *diag.Diagnostics, err error) {
	if _, ok := err.(api.NotFound); ok {
		state.RemoveResource(ctx)
	} else {
		errors.AddDiagError(diag, err)
	}
}
