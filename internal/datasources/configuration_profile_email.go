// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ datasource.DataSource = &configurationProfileEmailDataSource{}
)

func NewConfigurationProfileEmailDataSource() datasource.DataSource {
	return &configurationProfileEmailDataSource{}
}

type configurationProfileEmailDataSource struct {
	api *api.API
}

func (d *configurationProfileEmailDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_email"
}

func (d *configurationProfileEmailDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the email configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"from": schema.StringAttribute{
				Description: "The email from field value.",
				Computed:    true,
			},
			"smtp": schema.SingleNestedAttribute{
				Description: "SMTP configuration.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"server": schema.StringAttribute{
						Description: "SMTP server hostname or IP address.",
						Computed:    true,
					},
					"port": schema.StringAttribute{
						Description: "SMTP server port.",
						Computed:    true,
					},
					"ssl": schema.BoolAttribute{
						Description: "Whether SSL/TLS is enabled.",
						Computed:    true,
					},
					"username": schema.StringAttribute{
						Description: "Authentication username.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *configurationProfileEmailDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileEmailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileEmailDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.Read(ctx, api.ConfigurationProfileEmail)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(config.ID)
	data.Profile = types.StringValue(config.Profile)
	data.From = types.StringValue(config.Record.EmailConfiguration.FromEmail)

	if smtp := config.Record.EmailConfiguration.SMTP; smtp != nil {
		smtpAttrs := map[string]attr.Value{
			"server": types.StringValue(smtp.Server),
			"port":   types.StringValue(smtp.Port),
			"ssl":    types.BoolValue(smtp.SSL != nil && *smtp.SSL),
			"username": func() types.String {
				if smtp.Username != nil {
					return types.StringValue(*smtp.Username)
				}
				return types.StringNull()
			}(),
		}

		smtpObj, diags := types.ObjectValue(models.SMTP{}.AttributeTypes(), smtpAttrs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.SMTP = smtpObj
	} else {
		data.SMTP = types.ObjectNull(models.SMTP{}.AttributeTypes())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
