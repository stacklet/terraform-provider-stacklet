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
	_ resource.Resource                = &accountDiscoveryAWSResource{}
	_ resource.ResourceWithConfigure   = &accountDiscoveryAWSResource{}
	_ resource.ResourceWithImportState = &accountDiscoveryAWSResource{}
)

func NewAccountDiscoveryAWSResource() resource.Resource {
	return &accountDiscoveryAWSResource{}
}

type accountDiscoveryAWSResource struct {
	api *api.API
}

func (r *accountDiscoveryAWSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_discovery_aws"
}

func (r *accountDiscoveryAWSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an account discovery configuration for AWS.",
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
			"org_read_role": schema.StringAttribute{
				Description: "The ARN of an IAM role which has permission to read organization data.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"member_role": schema.StringAttribute{
				Description: "IAM role ARN template for AssetDB.",
				Required:    true,
			},
			"custodian_role": schema.StringAttribute{
				Description: "IAM role name or template for Cloud Custodian.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The organization ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *accountDiscoveryAWSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *accountDiscoveryAWSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.AccountDiscoveryAWSResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountDiscoveryAWSInput{
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueStringPointer(),
		CustodianRole: plan.CustodianRole.ValueString(),
		OrgReadRole:   plan.OrgReadRole.ValueString(),
		MemberRole:    plan.MemberRole.ValueString(),
	}
	accountDiscovery, err := r.api.AccountDiscovery.UpsertAWS(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	accountDiscovery, err = r.api.AccountDiscovery.UpdateSuspended(ctx, accountDiscovery.ID, plan.Suspended.ValueBool())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryAWSModel(&plan, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryAWSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountDiscoveryAWSResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountDiscovery, err := r.api.AccountDiscovery.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryAWSModel(&state, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountDiscoveryAWSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AccountDiscoveryAWSResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountDiscoveryAWSInput{
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueStringPointer(),
		CustodianRole: plan.CustodianRole.ValueString(),
		OrgReadRole:   plan.OrgReadRole.ValueString(),
		MemberRole:    plan.MemberRole.ValueString(),
	}
	accountDiscovery, err := r.api.AccountDiscovery.UpsertAWS(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	accountDiscovery, err = r.api.AccountDiscovery.UpdateSuspended(ctx, accountDiscovery.ID, plan.Suspended.ValueBool())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountDiscoveryAWSModel(&plan, accountDiscovery)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountDiscoveryAWSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Delete error", "Resource can't be deleted. Set `suspended = true` instead.")
}

func (r *accountDiscoveryAWSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

func updateAccountDiscoveryAWSModel(m *models.AccountDiscoveryAWSResource, accountDiscovery api.AccountDiscovery) {
	m.ID = types.StringValue(accountDiscovery.ID)
	m.Name = types.StringValue(accountDiscovery.Name)
	m.Description = tftypes.NullableString(accountDiscovery.Description)
	m.Suspended = types.BoolValue(accountDiscovery.Schedule.Suspended)
	m.OrgID = types.StringValue(accountDiscovery.Config.AWSConfig.OrgID)
	m.OrgReadRole = types.StringValue(accountDiscovery.Config.AWSConfig.OrgRole)
	m.CustodianRole = types.StringValue(accountDiscovery.Config.AWSConfig.CustodianRole)
}
