// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemavalidate"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &accountGroupResource{}
	_ resource.ResourceWithConfigure   = &accountGroupResource{}
	_ resource.ResourceWithImportState = &accountGroupResource{}
)

func NewAccountGroupResource() resource.Resource {
	return &accountGroupResource{}
}

type accountGroupResource struct {
	api *api.API
}

func (r *accountGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group"
}

func (r *accountGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an account group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the account group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the account group.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the account group.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account group (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
				Validators: []validator.String{
					schemavalidate.OneOfCloudProviders(),
				},
			},
			"regions": schema.ListAttribute{
				Description: "The list of regions for the account group (e.g., us-east-1, eu-west-2), for providers that require it.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *accountGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *accountGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.AccountGroupResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountGroupCreateInput{
		Name:        plan.Name.ValueString(),
		Provider:    plan.CloudProvider.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Regions:     api.StringsList(plan.Regions),
	}

	account_group, err := r.api.AccountGroup.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateAccountGroupModel(&plan, account_group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountGroupResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account_group, err := r.api.AccountGroup.Read(ctx, state.UUID.ValueString(), "")
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	r.updateAccountGroupModel(&state, account_group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AccountGroupResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountGroupUpdateInput{
		UUID:        plan.UUID.ValueString(),
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
		Regions:     api.StringsList(plan.Regions),
	}

	account_group, err := r.api.AccountGroup.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateAccountGroupModel(&plan, account_group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AccountGroupResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.AccountGroup.Delete(ctx, state.UUID.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *accountGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), req.ID)...)
}

func (r accountGroupResource) updateAccountGroupModel(m *models.AccountGroupResource, account_group *api.AccountGroup) {
	m.ID = types.StringValue(account_group.ID)
	m.UUID = types.StringValue(account_group.UUID)
	m.Name = types.StringValue(account_group.Name)
	m.Description = tftypes.NullableString(account_group.Description)
	m.CloudProvider = types.StringValue(account_group.Provider)
	m.Regions = tftypes.StringsList(account_group.Regions)
}
