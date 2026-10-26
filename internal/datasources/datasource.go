// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/api/platform"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

// apiDataSource is a datasource based on an API client.
type apiDataSource[T providerdata.APIClient] struct {
	api T
}

// Configure sets up API access for the datasource.
func (d *apiDataSource[T]) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if api, err := providerdata.GetAPI[T]("datasource", req.ProviderData); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else {
		d.api = api
	}
}

// platformAPIDataSource is a datasource based on the platform API.
type platformAPIDataSource = apiDataSource[*platform.API]

// iacAPIDataSource is a datasource based on the IaC API.
type iacAPIDataSource = apiDataSource[*iac.API] //nolint:unused
