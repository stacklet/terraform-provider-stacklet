package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ resource.Resource                = &accountDiscoveryResource{}
	_ resource.ResourceWithImportState = &accountDiscoveryResource{}
)

func NewAccountDiscoveryResource() resource.Resource {
	return &accountDiscoveryResource{}
}

type accountDiscoveryResource struct {
	client *graphql.Client
}

type accountDiscoveryResourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Enabled       types.Bool   `tfsdk:"enabled"`

	// AWS-specific fields
	OrgReadRole   types.String `tfsdk:"org_read_role"`
	MemberRole    types.String `tfsdk:"member_role"`
	CustodianRole types.String `tfsdk:"custodian_role"`

	// Azure-specific fields
	TenantID     types.String `tfsdk:"tenant_id"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`

	// GCP-specific fields
	OrgID            types.String `tfsdk:"org_id"`
	RootFolderIDs    types.List   `tfsdk:"root_folder_ids"`
	ExcludeFolderIDs types.List   `tfsdk:"exclude_folder_ids"`
	CredentialJSON   types.String `tfsdk:"credential_json"`

	// Common fields
	Suspended types.Bool   `tfsdk:"suspended"`
	Validity  types.Object `tfsdk:"validity"`
}

func (r *accountDiscoveryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_discovery"
}

func (r *accountDiscoveryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an account discovery configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account discovery configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the account discovery configuration.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Human-readable notes about the account discovery configuration.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Required:    true,
				Description: "The cloud provider for the account discovery (aws, azure, gcp, kubernetes, or tencentcloud).",
			},
			// AWS-specific fields
			"org_read_role": schema.StringAttribute{
				Description: "The ARN of an IAM role which has permission to read organization data (AWS only).",
				Optional:    true,
			},
			"member_role": schema.StringAttribute{
				Description: "Optional IAM role ARN template for AssetDB (AWS only).",
				Optional:    true,
			},
			"custodian_role": schema.StringAttribute{
				Description: "Optional IAM role name or ARN template for Cloud Custodian (AWS only).",
				Optional:    true,
			},
			// Azure-specific fields
			"tenant_id": schema.StringAttribute{
				Description: "The Azure tenant ID (Azure only).",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The Azure client ID (Azure only).",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The Azure client secret (Azure only).",
				Optional:    true,
				Sensitive:   true,
			},
			// GCP-specific fields
			"org_id": schema.StringAttribute{
				Description: "The GCP organization ID (GCP only).",
				Optional:    true,
			},
			"root_folder_ids": schema.ListAttribute{
				Description: "Optional list of GCP folder IDs to scan (GCP only).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"exclude_folder_ids": schema.ListAttribute{
				Description: "Optional list of GCP folder IDs to exclude from scanning (GCP only).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"credential_json": schema.StringAttribute{
				Description: "The contents of a JSON-formatted key file for a GCP service account (GCP only).",
				Optional:    true,
				Sensitive:   true,
			},
			// Common fields
			"suspended": schema.BoolAttribute{
				Description: "Whether account discovery is suspended.",
				Optional:    true,
			},
			"validity": schema.ObjectAttribute{
				Description: "Information about the most recent credential validation attempt.",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"valid":   types.BoolType,
					"message": types.StringType,
				},
			},
		},
	}
}

func (r *accountDiscoveryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *accountDiscoveryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accountDiscoveryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		UpsertAWSAccountDiscovery struct {
			AccountDiscovery struct {
				ID          string
				Name        string
				Description string
				Provider    string
				Config      struct {
					TypeName string `graphql:"__typename"`
					// AWS-specific fields
					OrgReadRole   string `graphql:"... on AWSAccountDiscoveryConfig"`
					MemberRole    string `graphql:"... on AWSAccountDiscoveryConfig"`
					CustodianRole string `graphql:"... on AWSAccountDiscoveryConfig"`
				}
				Schedule struct {
					Suspended bool
				}
				Validity struct {
					Valid   bool
					Message string
				}
			}
		} `graphql:"upsertAWSAccountDiscovery(input: $input)"`
		UpsertAzureAccountDiscovery struct {
			AccountDiscovery struct {
				ID          string
				Name        string
				Description string
				Provider    string
				Config      struct {
					TypeName string `graphql:"__typename"`
					// Azure-specific fields
					TenantID     string `graphql:"... on AzureAccountDiscoveryConfig"`
					ClientID     string `graphql:"... on AzureAccountDiscoveryConfig"`
					ClientSecret string `graphql:"... on AzureAccountDiscoveryConfig"`
				}
				Schedule struct {
					Suspended bool
				}
				Validity struct {
					Valid   bool
					Message string
				}
			}
		} `graphql:"upsertAzureAccountDiscovery(input: $input)"`
		UpsertGCPAccountDiscovery struct {
			AccountDiscovery struct {
				ID          string
				Name        string
				Description string
				Provider    string
				Config      struct {
					TypeName string `graphql:"__typename"`
					// GCP-specific fields
					OrgID            string   `graphql:"... on GCPAccountDiscoveryConfig"`
					RootFolderIDs    []string `graphql:"... on GCPAccountDiscoveryConfig"`
					ExcludeFolderIDs []string `graphql:"... on GCPAccountDiscoveryConfig"`
					CredentialJSON   string   `graphql:"... on GCPAccountDiscoveryConfig"`
				}
				Schedule struct {
					Suspended bool
				}
				Validity struct {
					Valid   bool
					Message string
				}
			}
		} `graphql:"upsertGCPAccountDiscovery(input: $input)"`
	}

	var err error
	var result struct {
		ID          string
		Name        string
		Description string
		Provider    string
		Config      any
		Schedule    struct {
			Suspended bool
		}
		Validity struct {
			Valid   bool
			Message string
		}
	}

	switch strings.ToLower(plan.CloudProvider.ValueString()) {
	case "aws":
		input := map[string]any{
			"name":        graphql.String(plan.Name.ValueString()),
			"orgReadRole": graphql.String(plan.OrgReadRole.ValueString()),
		}
		if !plan.Description.IsNull() {
			input["description"] = graphql.String(plan.Description.ValueString())
		}
		if !plan.MemberRole.IsNull() {
			input["memberRole"] = graphql.String(plan.MemberRole.ValueString())
		}
		if !plan.CustodianRole.IsNull() {
			input["custodianRole"] = graphql.String(plan.CustodianRole.ValueString())
		}
		variables := map[string]any{"input": input}
		err = r.client.Mutate(ctx, &mutation.UpsertAWSAccountDiscovery, variables)
		if err == nil {
			result.ID = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.ID
			result.Name = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.Name
			result.Description = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.Description
			result.Provider = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.Provider
			result.Schedule.Suspended = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.Schedule.Suspended
			result.Validity.Valid = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.Validity.Valid
			result.Validity.Message = mutation.UpsertAWSAccountDiscovery.AccountDiscovery.Validity.Message
		}

	case "azure":
		input := map[string]any{
			"name":         graphql.String(plan.Name.ValueString()),
			"tenantID":     graphql.String(plan.TenantID.ValueString()),
			"clientID":     graphql.String(plan.ClientID.ValueString()),
			"clientSecret": graphql.String(plan.ClientSecret.ValueString()),
		}
		if !plan.Description.IsNull() {
			input["description"] = graphql.String(plan.Description.ValueString())
		}
		variables := map[string]any{"input": input}
		err = r.client.Mutate(ctx, &mutation.UpsertAzureAccountDiscovery, variables)
		if err == nil {
			result.ID = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.ID
			result.Name = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.Name
			result.Description = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.Description
			result.Provider = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.Provider
			result.Schedule.Suspended = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.Schedule.Suspended
			result.Validity.Valid = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.Validity.Valid
			result.Validity.Message = mutation.UpsertAzureAccountDiscovery.AccountDiscovery.Validity.Message
		}

	case "gcp":
		input := map[string]any{
			"name":           graphql.String(plan.Name.ValueString()),
			"orgID":          graphql.String(plan.OrgID.ValueString()),
			"credentialJSON": graphql.String(plan.CredentialJSON.ValueString()),
		}
		if !plan.Description.IsNull() {
			input["description"] = graphql.String(plan.Description.ValueString())
		}
		if !plan.RootFolderIDs.IsNull() {
			var folderIDs []string
			plan.RootFolderIDs.ElementsAs(ctx, &folderIDs, false)
			input["rootFolderIDs"] = folderIDs
		}
		if !plan.ExcludeFolderIDs.IsNull() {
			var folderIDs []string
			plan.ExcludeFolderIDs.ElementsAs(ctx, &folderIDs, false)
			input["excludeFolderIDs"] = folderIDs
		}
		variables := map[string]any{"input": input}
		err = r.client.Mutate(ctx, &mutation.UpsertGCPAccountDiscovery, variables)
		if err == nil {
			result.ID = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.ID
			result.Name = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.Name
			result.Description = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.Description
			result.Provider = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.Provider
			result.Schedule.Suspended = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.Schedule.Suspended
			result.Validity.Valid = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.Validity.Valid
			result.Validity.Message = mutation.UpsertGCPAccountDiscovery.AccountDiscovery.Validity.Message
		}

	default:
		resp.Diagnostics.AddError(
			"Invalid Provider",
			fmt.Sprintf("Provider must be one of: aws, azure, gcp. Got: %s", plan.CloudProvider.ValueString()),
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account discovery, got error: %s", err))
		return
	}

	// Update schedule if suspended is set
	if !plan.Suspended.IsNull() {
		var scheduleMutation struct {
			UpdateAccountDiscoverySchedule struct {
				AccountDiscovery struct {
					Schedule struct {
						Suspended bool
					}
				}
			} `graphql:"updateAccountDiscoverySchedule(input: $input)"`
		}

		scheduleVariables := map[string]any{
			"input": map[string]any{
				"name":      graphql.String(result.Name),
				"suspended": graphql.Boolean(plan.Suspended.ValueBool()),
			},
		}

		err = r.client.Mutate(ctx, &scheduleMutation, scheduleVariables)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update account discovery schedule, got error: %s", err))
			return
		}

		result.Schedule.Suspended = scheduleMutation.UpdateAccountDiscoverySchedule.AccountDiscovery.Schedule.Suspended
	}

	plan.ID = types.StringValue(result.ID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.CloudProvider = types.StringValue(result.Provider)
	plan.Suspended = types.BoolValue(result.Schedule.Suspended)

	validityObj, diags := types.ObjectValue(
		map[string]attr.Type{
			"valid":   types.BoolType,
			"message": types.StringType,
		},
		map[string]attr.Value{
			"valid":   types.BoolValue(result.Validity.Valid),
			"message": types.StringValue(result.Validity.Message),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Validity = validityObj

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountDiscoveryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query struct {
		AccountDiscovery struct {
			ID          string
			Name        string
			Description string
			Provider    string
			Config      struct {
				TypeName string `graphql:"__typename"`
				// AWS-specific fields
				OrgReadRole   string `graphql:"... on AWSAccountDiscoveryConfig"`
				MemberRole    string `graphql:"... on AWSAccountDiscoveryConfig"`
				CustodianRole string `graphql:"... on AWSAccountDiscoveryConfig"`
				// Azure-specific fields
				TenantID     string `graphql:"... on AzureAccountDiscoveryConfig"`
				ClientID     string `graphql:"... on AzureAccountDiscoveryConfig"`
				ClientSecret string `graphql:"... on AzureAccountDiscoveryConfig"`
				// GCP-specific fields
				OrgID            string   `graphql:"... on GCPAccountDiscoveryConfig"`
				RootFolderIDs    []string `graphql:"... on GCPAccountDiscoveryConfig"`
				ExcludeFolderIDs []string `graphql:"... on GCPAccountDiscoveryConfig"`
				CredentialJSON   string   `graphql:"... on GCPAccountDiscoveryConfig"`
			}
			Schedule struct {
				Suspended bool
			}
			Validity struct {
				Valid   bool
				Message string
			}
		} `graphql:"accountDiscovery(name: $name)"`
	}

	variables := map[string]any{
		"name": graphql.String(state.Name.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account discovery, got error: %s", err))
		return
	}

	if query.AccountDiscovery.Name == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(query.AccountDiscovery.ID)
	state.Name = types.StringValue(query.AccountDiscovery.Name)
	state.Description = types.StringValue(query.AccountDiscovery.Description)
	state.CloudProvider = types.StringValue(query.AccountDiscovery.Provider)
	state.Suspended = types.BoolValue(query.AccountDiscovery.Schedule.Suspended)

	// Set provider-specific fields
	switch query.AccountDiscovery.Config.TypeName {
	case "AWSAccountDiscoveryConfig":
		state.OrgReadRole = types.StringValue(query.AccountDiscovery.Config.OrgReadRole)
		state.MemberRole = types.StringValue(query.AccountDiscovery.Config.MemberRole)
		state.CustodianRole = types.StringValue(query.AccountDiscovery.Config.CustodianRole)
	case "AzureAccountDiscoveryConfig":
		state.TenantID = types.StringValue(query.AccountDiscovery.Config.TenantID)
		state.ClientID = types.StringValue(query.AccountDiscovery.Config.ClientID)
		state.ClientSecret = types.StringValue(query.AccountDiscovery.Config.ClientSecret)
	case "GCPAccountDiscoveryConfig":
		state.OrgID = types.StringValue(query.AccountDiscovery.Config.OrgID)
		rootFolderIDs, diags := types.ListValueFrom(ctx, types.StringType, query.AccountDiscovery.Config.RootFolderIDs)
		resp.Diagnostics.Append(diags...)
		state.RootFolderIDs = rootFolderIDs
		excludeFolderIDs, diags := types.ListValueFrom(ctx, types.StringType, query.AccountDiscovery.Config.ExcludeFolderIDs)
		resp.Diagnostics.Append(diags...)
		state.ExcludeFolderIDs = excludeFolderIDs
		state.CredentialJSON = types.StringValue(query.AccountDiscovery.Config.CredentialJSON)
	}

	validityObj, diags := types.ObjectValue(
		map[string]attr.Type{
			"valid":   types.BoolType,
			"message": types.StringType,
		},
		map[string]attr.Value{
			"valid":   types.BoolValue(query.AccountDiscovery.Validity.Valid),
			"message": types.StringValue(query.AccountDiscovery.Validity.Message),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Validity = validityObj

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountDiscoveryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accountDiscoveryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The Update implementation is similar to Create, using the same mutations
	// Just call Create since it uses upsert mutations
	r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, (*resource.CreateResponse)(resp))
}

func (r *accountDiscoveryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accountDiscoveryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// There's no explicit delete mutation in the GraphQL schema
	// We'll suspend the account discovery instead
	var mutation struct {
		UpdateAccountDiscoverySchedule struct {
			AccountDiscovery struct {
				Name string
			}
		} `graphql:"updateAccountDiscoverySchedule(input: $input)"`
	}

	variables := map[string]any{
		"input": map[string]any{
			"name":      graphql.String(state.Name.ValueString()),
			"suspended": graphql.Boolean(true),
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to suspend account discovery, got error: %s", err))
		return
	}
}

func (r *accountDiscoveryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is the name of the account discovery configuration
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}
