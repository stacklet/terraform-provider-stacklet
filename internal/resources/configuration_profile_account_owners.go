// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &configurationProfileAccountOwnersResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileAccountOwnersResource{}
	_ resource.ResourceWithImportState = &configurationProfileAccountOwnersResource{}
)

func NewConfigurationProfileAccountOwnersResource() resource.Resource {
	return &configurationProfileAccountOwnersResource{}
}

type configurationProfileAccountOwnersResource struct {
	api *api.API
}

func (r *configurationProfileAccountOwnersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_account_owners"
}

func (r *configurationProfileAccountOwnersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the account owners configuration profile.

The profile is global, adding multiple resources of this kind will cause them to override each other.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default": schema.ListNestedAttribute{
				Description: "List of default account owners.",
				Optional:    true,
				Computed:    true,
				Default:     tftypes.EmptyListDefault(r.ownersListType()),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account": schema.StringAttribute{
							Description: "The account identifier.",
							Required:    true,
						},
						"owners": schema.ListAttribute{
							Description: "List of owner addresses for this account.",
							Required:    true,
							ElementType: types.StringType,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
			"org_domain": schema.StringAttribute{
				Description: "The organization domain to append to users for matching.",
				Optional:    true,
			},
			"org_domain_tag": schema.StringAttribute{
				Description: "The name of the tag to look up the organization domain.",
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				Description: "List of tag names to look up the resource owner.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     tftypes.EmptyListDefault(types.StringType),
			},
		},
	}
}

func (r *configurationProfileAccountOwnersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileAccountOwnersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileAccountOwnersResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileConfig, err := r.api.ConfigurationProfile.ReadAccountOwners(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(*profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileAccountOwnersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ConfigurationProfileAccountOwnersResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultOwners, diags := r.getDefaultOwners(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountOwnersConfiguration{
		Default:      defaultOwners,
		OrgDomain:    plan.OrgDomain.ValueStringPointer(),
		OrgDomainTag: plan.OrgDomainTag.ValueStringPointer(),
		Tags:         api.StringsList(plan.Tags),
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertAccountOwners(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(*profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileAccountOwnersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ConfigurationProfileAccountOwnersResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultOwners, diags := r.getDefaultOwners(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountOwnersConfiguration{
		Default:      defaultOwners,
		OrgDomain:    plan.OrgDomain.ValueStringPointer(),
		OrgDomainTag: plan.OrgDomainTag.ValueStringPointer(),
		Tags:         api.StringsList(plan.Tags),
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertAccountOwners(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(*profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileAccountOwnersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileAccountOwnersResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteAccountOwners(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileAccountOwnersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileAccountOwners))...)
}

func (r configurationProfileAccountOwnersResource) getDefaultOwners(ctx context.Context, m models.ConfigurationProfileAccountOwnersResource) ([]api.AccountOwners, diag.Diagnostics) {
	if m.Default.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	owners := []api.AccountOwners{}
	for i, elem := range m.Default.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("default.%d", i)),
				"Invalid account owners default settings",
				"Account owners default block is invalid.",
			)
			return nil, diags
		}
		var o models.AccountOwners
		if diags := block.As(ctx, &o, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}

		owners = append(
			owners,
			api.AccountOwners{
				Account: o.Account.ValueString(),
				Owners:  api.StringsList(o.Owners),
			},
		)
	}
	return owners, diags
}

func (r configurationProfileAccountOwnersResource) ownersListType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"account": types.StringType,
			"owners":  types.ListType{ElemType: types.StringType},
		},
	}
}
