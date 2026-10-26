// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/api/platform"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

// apiResource is a resource based on an API client.
type apiResource[T providerdata.APIClient] struct {
	api T
}

// Configure sets up API access for the resource.
func (r *apiResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if api, err := providerdata.GetAPI[T]("resource", req.ProviderData); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else {
		r.api = api
	}
}

// platformAPIResource is a resource based on the platform API.
type platformAPIResource = apiResource[*platform.API]

// iacAPIResource is a resource based on the IaC API.
type iacAPIResource = apiResource[*iac.API] //nolint:unused
