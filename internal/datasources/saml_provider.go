// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var _ datasource.DataSource = &samlProviderDataSource{}

type samlProviderDataSource struct {
	apiDataSource
}

func (d *samlProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saml_provider"
}

func (d *samlProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a SAML identity provider by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the SAML provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The unique name identifying the provider.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the SAML provider.",
				Computed:    true,
			},
			"metadata_url": schema.StringAttribute{
				Description: "URL of the SAML metadata endpoint, if the provider is configured with one.",
				Computed:    true,
			},
			"metadata_xml": schema.StringAttribute{
				Description: "Inline SAML metadata XML document, if the provider is configured with one.",
				Computed:    true,
			},
			"enable_signout": schema.BoolAttribute{
				Description: "Whether the identity provider signout flow is enabled for this provider.",
				Computed:    true,
			},
			"idp_alias": schema.StringAttribute{
				Description: "A unique human-facing alias for the provider.",
				Computed:    true,
			},
		},
	}
}

func (d *samlProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.SAMLProviderDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	samlProvider, err := d.api.SAMLProvider.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(samlProvider)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
