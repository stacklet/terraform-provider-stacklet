// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ datasource.DataSource = &configurationProfileJiraDataSource{}
)

func NewConfigurationProfileJiraDataSource() datasource.DataSource {
	return &configurationProfileJiraDataSource{}
}

type configurationProfileJiraDataSource struct {
	api *api.API
}

func (d *configurationProfileJiraDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_jira"
}

func (d *configurationProfileJiraDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Jira configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The Jira instance URL.",
				Computed:    true,
			},
			"user": schema.StringAttribute{
				Description: "The Jira instance authentication username.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"project": schema.ListNestedBlock{
				Description: "Jira project configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"closed_status": schema.StringAttribute{
							Description: "The state for closed tickets.",
							Computed:    true,
						},
						"issue_type": schema.StringAttribute{
							Description: "The type of issue to use for tickets.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the project.",
							Computed:    true,
						},
						"project": schema.StringAttribute{
							Description: "The ID of the project.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *configurationProfileJiraDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileJiraDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileJiraDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadJira(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(config.ID)
	data.Profile = types.StringValue(config.Profile)
	data.URL = types.StringPointerValue(config.Record.JiraConfiguration.URL)
	data.User = types.StringValue(config.Record.JiraConfiguration.User)

	updater := modelupdate.NewConfigurationProfileUpdater(*config)
	projects, diags := updater.JiraProjects()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Projects = projects

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
