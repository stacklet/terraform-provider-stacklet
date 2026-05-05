// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &gcpIntegrationResource{}
	_ resource.ResourceWithConfigure   = &gcpIntegrationResource{}
	_ resource.ResourceWithImportState = &gcpIntegrationResource{}
)

func newGCPIntegrationResource() resource.Resource {
	return &gcpIntegrationResource{}
}

type gcpIntegrationResource struct {
	apiResource
}

func (r *gcpIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_integration"
}

func (r *gcpIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	terraformModuleAttrs := map[string]schema.Attribute{
		"repository_url": schema.StringAttribute{
			Description: "The Terraform module repository URL.",
			Computed:    true,
		},
		"source": schema.StringAttribute{
			Description: "The Terraform module source.",
			Computed:    true,
		},
		"version": schema.StringAttribute{
			Description: "The Terraform module version.",
			Computed:    true,
		},
		"variables_json": schema.StringAttribute{
			Description: "The Terraform module variables as JSON.",
			Computed:    true,
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages a GCP integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the GCP integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The key identifying the GCP integration.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"customer_config_input": schema.SingleNestedAttribute{
				Description: "Customer-provided configuration defining the desired scope and surfaces Terraform to deploy it.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"infrastructure": schema.SingleNestedAttribute{
						Description: "The GCP project where Stacklet resources will be deployed.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"project_id": schema.StringAttribute{
								Description: "The GCP project ID for Stacklet infrastructure.",
								Required:    true,
							},
							"resource_location": schema.StringAttribute{
								Description: "The GCP location for Stacklet resources.",
								Required:    true,
							},
							"resource_prefix": schema.StringAttribute{
								Description: "The prefix prepended to all Stacklet-managed resource names.",
								Optional:    true,
							},
							"create_project": schema.SingleNestedAttribute{
								Description: "Configuration for creating the Stacklet infrastructure project, if applicable.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"billing_account_id": schema.StringAttribute{
										Description: "The billing account to associate with the created project.",
										Required:    true,
									},
									"org_id": schema.StringAttribute{
										Description: "The organization in which to create the project.",
										Optional:    true,
									},
									"folder_id": schema.StringAttribute{
										Description: "The folder in which to create the project.",
										Optional:    true,
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
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"org_id": schema.StringAttribute{
									Description: "The GCP organization ID.",
									Required:    true,
								},
								"folder_ids": schema.ListAttribute{
									Description: "Folders to manage within the organization.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"project_ids": schema.ListAttribute{
									Description: "Individual projects to manage within the organization.",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"cost_sources": schema.ListNestedAttribute{
						Description: "Billing export tables to use for cost data.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"billing_table": schema.StringAttribute{
									Description: "The BigQuery table containing billing export data.",
									Required:    true,
								},
							},
						},
					},
					"security_contexts": schema.ListNestedAttribute{
						Description: "Named sets of extra roles to grant beyond the integration defaults.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The security context name.",
									Required:    true,
								},
								"extra_roles": schema.ListAttribute{
									Description: "Additional GCP roles to grant for this context.",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"access_config_blob_input": schema.StringAttribute{
				Description: "Base64-encoded JSON access config as output by the Terraform module.",
				Optional:    true,
			},
			"customer_config": schema.SingleNestedAttribute{
				Description: "Customer-provided configuration returned by the API.",
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
										Computed:    true,
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
					"terraform_module": schema.SingleNestedAttribute{
						Description: "Terraform module to deploy this integration, computed by Stacklet.",
						Computed:    true,
						Attributes:  terraformModuleAttrs,
					},
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

func (r *gcpIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.GCPIntegrationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, d := buildGCPIntegrationInput(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.api.GCPIntegration.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *gcpIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.GCPIntegrationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.api.GCPIntegration.Read(ctx, state.Key.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *gcpIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.GCPIntegrationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, d := buildGCPIntegrationInput(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.api.GCPIntegration.Upsert(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *gcpIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.GCPIntegrationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.GCPIntegration.Delete(ctx, state.Key.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	}
}

func (r *gcpIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"key"})
}

func buildGCPIntegrationInput(ctx context.Context, plan models.GCPIntegrationResource) (api.GCPIntegrationInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	customerConfig, d := buildCustomerConfigInput(ctx, plan.CustomerConfigInput)
	diags.Append(d...)

	input := api.GCPIntegrationInput{
		Key:              plan.Key.ValueString(),
		AccessConfigBlob: plan.AccessConfigBlobInput.ValueStringPointer(),
		CustomerConfig:   customerConfig,
	}
	return input, diags
}

func buildCustomerConfigInput(ctx context.Context, obj types.Object) (*api.GCPIntegrationCustomerConfigInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var m models.GCPIntegrationCustomerConfigInputModel
	diags.Append(obj.As(ctx, &m, objectAsOptions)...)
	if diags.HasError() {
		return nil, diags
	}

	infrastructure, d := buildInfrastructureInput(ctx, m.Infrastructure)
	diags.Append(d...)
	organizations, d := buildOrganizationsInput(ctx, m.Organizations)
	diags.Append(d...)
	costSources, d := buildCostSourcesInput(ctx, m.CostSources)
	diags.Append(d...)
	securityContexts, d := buildSecurityContextsInput(ctx, m.SecurityContexts)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.GCPIntegrationCustomerConfigInput{
		Infrastructure:   infrastructure,
		Organizations:    organizations,
		CostSources:      costSources,
		SecurityContexts: securityContexts,
	}, diags
}

func buildInfrastructureInput(ctx context.Context, obj types.Object) (*api.GCPIntegrationCustomerInfrastructureInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var m models.GCPIntegrationCustomerInfraModel
	diags.Append(obj.As(ctx, &m, objectAsOptions)...)
	if diags.HasError() {
		return nil, diags
	}

	createProject, d := buildCreateProjectInput(ctx, m.CreateProject)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.GCPIntegrationCustomerInfrastructureInput{
		ProjectID:        m.ProjectID.ValueStringPointer(),
		ResourceLocation: m.ResourceLocation.ValueStringPointer(),
		ResourcePrefix:   m.ResourcePrefix.ValueStringPointer(),
		CreateProject:    createProject,
	}, diags
}

func buildCreateProjectInput(ctx context.Context, obj types.Object) (*api.GCPIntegrationCustomerCreateProjectInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var m models.GCPIntegrationCustomerCreateProjectModel
	diags.Append(obj.As(ctx, &m, objectAsOptions)...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.GCPIntegrationCustomerCreateProjectInput{
		BillingAccountID: m.BillingAccountID.ValueStringPointer(),
		OrgID:            m.OrgID.ValueStringPointer(),
		FolderID:         m.FolderID.ValueStringPointer(),
		Labels:           api.NewTagsList(m.Labels),
	}, diags
}

func buildOrganizationsInput(ctx context.Context, l types.List) (*[]api.GCPIntegrationCustomerOrganizationInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if l.IsNull() || l.IsUnknown() {
		return nil, diags
	}

	var ms []models.GCPIntegrationCustomerOrgModel
	diags.Append(l.ElementsAs(ctx, &ms, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]api.GCPIntegrationCustomerOrganizationInput, len(ms))
	for i, m := range ms {
		result[i] = api.GCPIntegrationCustomerOrganizationInput{
			OrgID:      m.OrgID.ValueString(),
			FolderIDs:  api.StringsList(m.FolderIDs),
			ProjectIDs: api.StringsList(m.ProjectIDs),
		}
	}
	return &result, diags
}

func buildCostSourcesInput(ctx context.Context, l types.List) (*[]api.GCPIntegrationCustomerCostSourceInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if l.IsNull() || l.IsUnknown() {
		return nil, diags
	}

	var ms []models.GCPIntegrationCustomerCostSourceModel
	diags.Append(l.ElementsAs(ctx, &ms, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]api.GCPIntegrationCustomerCostSourceInput, len(ms))
	for i, m := range ms {
		result[i] = api.GCPIntegrationCustomerCostSourceInput{
			BillingTable: m.BillingTable.ValueString(),
		}
	}
	return &result, diags
}

func buildSecurityContextsInput(ctx context.Context, l types.List) (*[]api.GCPIntegrationCustomerSecurityContextInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if l.IsNull() || l.IsUnknown() {
		return nil, diags
	}

	var ms []models.GCPIntegrationCustomerSecurityContextModel
	diags.Append(l.ElementsAs(ctx, &ms, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]api.GCPIntegrationCustomerSecurityContextInput, len(ms))
	for i, m := range ms {
		result[i] = api.GCPIntegrationCustomerSecurityContextInput{
			Name:       m.Name.ValueString(),
			ExtraRoles: api.StringsList(m.ExtraRoles),
		}
	}
	return &result, diags
}
