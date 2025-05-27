package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &accountDiscoveryGCPResource{}
	_ resource.ResourceWithConfigure   = &accountDiscoveryGCPResource{}
	_ resource.ResourceWithImportState = &accountDiscoveryGCPResource{}
)

func NewAccountDiscoveryGCPResource() resource.Resource {
	return &accountDiscoveryGCPResource{}
}

type accountDiscoveryGCPResource struct {
	api *api.API
}

func (r *accountDiscoveryGCPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_discovery_gcp"
}

func (r *accountDiscoveryGCPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an account discovery configuration for GCP.",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Human-readable notes about the account discovery configuration.",
				Optional:    true,
			},
			"suspended": schema.BoolAttribute{
				Description: "Whether the discovery schedule is suspended.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"root_folder_ids": schema.ListAttribute{
				Description: "List of GCP folder IDs to scan.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"exclude_folder_ids": schema.ListAttribute{
				Description: "List of GCP folder IDs to exclude from scanning.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"credential_json_wo": schema.StringAttribute{
				Description: "The contents of a JSON-formatted key file for a GCP service account.",
				Required:    true,
				WriteOnly:   true,
			},
			"credential_json_wo_version": schema.StringAttribute{
				Description: "The version for key GCP service account key file. Must be changed to update credential_json_wo.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The GCP organization ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_email": schema.StringAttribute{
				Description: "The client email for the configuration.",
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID for the configuration.",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID for the configuration.",
				Computed:    true,
			},
			"private_key_id": schema.StringAttribute{
				Description: "The private key ID.",
				Computed:    true,
			},
		},
	}
}

func (r *accountDiscoveryGCPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *accountDiscoveryGCPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.AccountDiscoveryGCPResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountDiscoveryGCPInput{
		Name:             plan.Name.ValueString(),
		Description:      plan.Description.ValueStringPointer(),
		OrgID:            plan.OrgID.ValueString(),
		RootFolderIDs:    api.StringsList(plan.RootFolderIDs),
		ExcludeFolderIDs: api.StringsList(plan.ExcludeFolderIDs),
		CredentialJSON:   config.CredentialJSON.ValueStringPointer(),
	}
	accountDiscovery, err := r.api.AccountDiscovery.UpsertGCP(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	accountDiscovery, err = r.api.AccountDiscovery.UpdateSuspended(ctx, accountDiscovery.ID, plan.Suspended.ValueBool())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryGCPModel(&plan, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryGCPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountDiscoveryGCPResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountDiscovery, err := r.api.AccountDiscovery.Read(ctx, state.Name.ValueString())
	if err != nil {
		errors.HandleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryGCPModel(&state, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountDiscoveryGCPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, config, state models.AccountDiscoveryGCPResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var credentialJSON *string
	if state.CredentialJSONVersion != plan.CredentialJSONVersion {
		credentialJSON = config.CredentialJSON.ValueStringPointer()
	}
	input := api.AccountDiscoveryGCPInput{
		Name:             plan.Name.ValueString(),
		Description:      plan.Description.ValueStringPointer(),
		OrgID:            plan.OrgID.ValueString(),
		RootFolderIDs:    api.StringsList(plan.RootFolderIDs),
		ExcludeFolderIDs: api.StringsList(plan.ExcludeFolderIDs),
		CredentialJSON:   credentialJSON, // this might be nil, in which case it's not updated
	}
	accountDiscovery, err := r.api.AccountDiscovery.UpsertGCP(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	accountDiscovery, err = r.api.AccountDiscovery.UpdateSuspended(ctx, accountDiscovery.ID, plan.Suspended.ValueBool())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryGCPModel(&plan, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryGCPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Delete error", "Resource can't be deleted. Set `suspended = true` instead.")
}

func (r *accountDiscoveryGCPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

func updateAccountDiscoveryGCPModel(m *models.AccountDiscoveryGCPResource, accountDiscovery api.AccountDiscovery) {
	m.ID = types.StringValue(accountDiscovery.ID)
	m.Name = types.StringValue(accountDiscovery.Name)
	m.Description = tftypes.NullableString(accountDiscovery.Description)
	m.Suspended = types.BoolValue(accountDiscovery.Schedule.Suspended)
	m.ClientEmail = types.StringValue(accountDiscovery.Config.GCPConfig.ClientEmail)
	m.ClientID = types.StringValue(accountDiscovery.Config.GCPConfig.ClientID)
	m.OrgID = types.StringValue(accountDiscovery.Config.GCPConfig.OrgID)
	m.RootFolderIDs = tftypes.StringsList(accountDiscovery.Config.GCPConfig.RootFolderIDs)
	m.ExcludeFolderIDs = tftypes.StringsList(accountDiscovery.Config.GCPConfig.ExcludeFolderIDs)
	m.ProjectID = types.StringValue(accountDiscovery.Config.GCPConfig.ProjectID)
	m.PrivateKeyID = types.StringValue(accountDiscovery.Config.GCPConfig.PrivateKeyID)
}
