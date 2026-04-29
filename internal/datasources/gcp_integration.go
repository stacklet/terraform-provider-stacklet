// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ datasource.DataSource = &gcpIntegrationDataSource{}
)

func newGCPIntegrationDataSource() datasource.DataSource {
	return &gcpIntegrationDataSource{}
}

type gcpIntegrationDataSource struct {
	apiDataSource
}

func (d *gcpIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_integration"
}

func (d *gcpIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a GCP integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the GCP integration.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "The key identifying the GCP integration.",
				Required:    true,
			},
			"customer_config": schema.SingleNestedAttribute{
				Description: "Customer-provided configuration defining the desired scope and surfaces Terraform to deploy it.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"infrastructure": schema.SingleNestedAttribute{
						Description: "The GCP project where Stacklet resources will be deployed.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"project_id": schema.StringAttribute{
								Description: "The GCP project ID for Stacklet infrastructure.",
								Computed:    true,
							},
							"resource_location": schema.StringAttribute{
								Description: "The GCP location for Stacklet resources.",
								Computed:    true,
							},
							"resource_prefix": schema.StringAttribute{
								Description: "The prefix prepended to all Stacklet-managed resource names.",
								Computed:    true,
							},
							"create_project": schema.SingleNestedAttribute{
								Description: "Configuration for creating the Stacklet infrastructure project, if applicable.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"billing_account_id": schema.StringAttribute{
										Description: "The billing account to associate with the created project.",
										Computed:    true,
									},
									"org_id": schema.StringAttribute{
										Description: "The organization in which to create the project.",
										Computed:    true,
									},
									"folder_id": schema.StringAttribute{
										Description: "The folder in which to create the project.",
										Computed:    true,
									},
									"labels": schema.MapAttribute{
										Description: "Labels applied to the created project.",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
					"organizations": schema.ListNestedAttribute{
						Description: "GCP organizations, folders, and projects Stacklet will manage.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"org_id": schema.StringAttribute{
									Description: "The GCP organization ID.",
									Computed:    true,
								},
								"folder_ids": schema.ListAttribute{
									Description: "Folders to manage within the organization.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"project_ids": schema.ListAttribute{
									Description: "Individual projects to manage within the organization.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"cost_sources": schema.ListNestedAttribute{
						Description: "Billing export tables to use for cost data.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"billing_table": schema.StringAttribute{
									Description: "The BigQuery table containing billing export data.",
									Computed:    true,
								},
							},
						},
					},
					"security_contexts": schema.ListNestedAttribute{
						Description: "Named sets of extra roles to grant beyond the integration defaults.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The security context name.",
									Computed:    true,
								},
								"extra_roles": schema.ListAttribute{
									Description: "Additional GCP roles to grant for this context.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"terraform_module": models.TerraformModule{}.DataSourceSchemaAttribute(),
				},
			},
			"access_config": schema.SingleNestedAttribute{
				Description: "Live connection details populated after Terraform has been applied; null until then.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"infrastructure": schema.SingleNestedAttribute{
						Description: "The GCP infrastructure project and its relay and WIF configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"project_id": schema.StringAttribute{
								Description: "The GCP project hosting Stacklet infrastructure.",
								Computed:    true,
							},
							"relay": schema.SingleNestedAttribute{
								Description: "Event relay identity credential.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"oauth_id": schema.StringAttribute{
										Description: "The numeric unique_id of the events-relay Google Service Account.",
										Computed:    true,
									},
								},
							},
							"wif": schema.SingleNestedAttribute{
								Description: "Workload Identity Federation configuration.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"audience": schema.StringAttribute{
										Description: "The full resource name URI of the GCP WIF pool provider.",
										Computed:    true,
									},
									"principals": schema.SingleNestedAttribute{
										Description: "WIF service account principals by role.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"read_only": schema.StringAttribute{
												Description: "Service account email for read-only resource sync.",
												Computed:    true,
											},
											"cost_query": schema.StringAttribute{
												Description: "Service account email for cost queries.",
												Computed:    true,
											},
										},
									},
								},
							},
							"baseline_roles": schema.ListAttribute{
								Description: "Roles applied to all security contexts.",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
					"organizations": schema.ListNestedAttribute{
						Description: "Accessible GCP organizations.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "The GCP organization ID.",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "The GCP organization name.",
									Computed:    true,
								},
								"folders": schema.ListNestedAttribute{
									Description: "Accessible folders within the organization.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed: true,
											},
											"name": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
								"projects": schema.ListNestedAttribute{
									Description: "Accessible projects within the organization.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed: true,
											},
											"number": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"cost_sources": schema.ListNestedAttribute{
						Description: "Billing export tables and their locations.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"billing_table": schema.StringAttribute{
									Description: "The BigQuery table containing billing export data.",
									Computed:    true,
								},
								"location": schema.StringAttribute{
									Description: "The BigQuery dataset location.",
									Computed:    true,
								},
							},
						},
					},
					"security_contexts": schema.ListNestedAttribute{
						Description: "Named sets of roles granted to principals.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The security context name.",
									Computed:    true,
								},
								"extra_roles": schema.ListAttribute{
									Description: "GCP roles granted in addition to the baseline ones.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"principal": schema.StringAttribute{
									Description: "The GCP identity granted these roles.",
									Computed:    true,
								},
							},
						},
					},
					"roundtrip_digest": schema.StringAttribute{
						Description: "Digest of the customer config at deployment time.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *gcpIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.GCPIntegrationDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := d.api.GCPIntegration.Read(ctx, data.Key.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
