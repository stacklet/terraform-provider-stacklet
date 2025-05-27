// Copyright (c) 2025 - Stacklet, Inc.

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
	_ resource.Resource                = &accountDiscoveryAzureResource{}
	_ resource.ResourceWithConfigure   = &accountDiscoveryAzureResource{}
	_ resource.ResourceWithImportState = &accountDiscoveryAzureResource{}
)

func NewAccountDiscoveryAzureResource() resource.Resource {
	return &accountDiscoveryAzureResource{}
}

type accountDiscoveryAzureResource struct {
	api *api.API
}

func (r *accountDiscoveryAzureResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_discovery_azure"
}

func (r *accountDiscoveryAzureResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an account discovery configuration for Azure.",
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
			"tenant_id": schema.StringAttribute{
				Description: "The Azure tenant ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The Azure client ID.",
				Required:    true,
			},
			"client_secret_wo": schema.StringAttribute{
				Description: "The Azure client secret. This is not stored in state and is only updated when client_id is changed.",
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
		},
	}
}

func (r *accountDiscoveryAzureResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *accountDiscoveryAzureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.AccountDiscoveryAzureResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountDiscoveryAzureInput{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueStringPointer(),
		ClientID:     plan.ClientID.ValueStringPointer(),
		ClientSecret: config.ClientSecret.ValueStringPointer(),
		TenantID:     plan.TenantID.ValueString(),
	}
	accountDiscovery, err := r.api.AccountDiscovery.UpsertAzure(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	accountDiscovery, err = r.api.AccountDiscovery.UpdateSuspended(ctx, accountDiscovery.ID, plan.Suspended.ValueBool())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryAzureModel(&plan, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryAzureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountDiscoveryAzureResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountDiscovery, err := r.api.AccountDiscovery.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryAzureModel(&state, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountDiscoveryAzureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, config, state models.AccountDiscoveryAzureResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clientID, clientSecret *string
	// ID and secret must be set together, if not they're both left nil in the API call
	if state.ClientID != plan.ClientID {
		clientID = plan.ClientID.ValueStringPointer()
		clientSecret = config.ClientSecret.ValueStringPointer()
	}
	input := api.AccountDiscoveryAzureInput{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueStringPointer(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     plan.TenantID.ValueString(),
	}
	accountDiscovery, err := r.api.AccountDiscovery.UpsertAzure(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	accountDiscovery, err = r.api.AccountDiscovery.UpdateSuspended(ctx, accountDiscovery.ID, plan.Suspended.ValueBool())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryAzureModel(&plan, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryAzureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Delete error", "Resource can't be deleted. Set `suspended = true` instead.")
}

func (r *accountDiscoveryAzureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

func updateAccountDiscoveryAzureModel(m *models.AccountDiscoveryAzureResource, accountDiscovery api.AccountDiscovery) {
	m.ID = types.StringValue(accountDiscovery.ID)
	m.Name = types.StringValue(accountDiscovery.Name)
	m.Description = tftypes.NullableString(accountDiscovery.Description)
	m.Suspended = types.BoolValue(accountDiscovery.Schedule.Suspended)
	m.ClientID = types.StringValue(accountDiscovery.Config.AzureConfig.ClientID)
	m.TenantID = types.StringValue(accountDiscovery.Config.AzureConfig.TenantID)
}
