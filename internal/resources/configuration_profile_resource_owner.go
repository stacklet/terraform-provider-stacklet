// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &configurationProfileResourceOwnerResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileResourceOwnerResource{}
	_ resource.ResourceWithImportState = &configurationProfileResourceOwnerResource{}
)

func NewConfigurationProfileResourceOwnerResource() resource.Resource {
	return &configurationProfileResourceOwnerResource{}
}

type configurationProfileResourceOwnerResource struct {
	api *api.API
}

func (r *configurationProfileResourceOwnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_resource_owner"
}

func (r *configurationProfileResourceOwnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the resource owner configuration profile.

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
			"default": schema.ListAttribute{
				Description: "List of default resource owners.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     tftypes.DefaultStringListEmpty(),
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
				Default:     tftypes.DefaultStringListEmpty(),
			},
		},
	}
}

func (r *configurationProfileResourceOwnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileResourceOwnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileResourceOwnerResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.api.ConfigurationProfile.ReadResourceOwner(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	r.updateResourceOwnerModel(&state, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileResourceOwnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ConfigurationProfileResourceOwnerResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.ResourceOwnerConfiguration{
		Default:      api.StringsList(plan.Default),
		OrgDomain:    plan.OrgDomain.ValueStringPointer(),
		OrgDomainTag: plan.OrgDomainTag.ValueStringPointer(),
		Tags:         api.StringsList(plan.Tags),
	}
	config, err := r.api.ConfigurationProfile.UpsertResourceOwner(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateResourceOwnerModel(&plan, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileResourceOwnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ConfigurationProfileResourceOwnerResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.ResourceOwnerConfiguration{
		Default:      api.StringsList(plan.Default),
		OrgDomain:    plan.OrgDomain.ValueStringPointer(),
		OrgDomainTag: plan.OrgDomainTag.ValueStringPointer(),
		Tags:         api.StringsList(plan.Tags),
	}
	config, err := r.api.ConfigurationProfile.UpsertResourceOwner(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	r.updateResourceOwnerModel(&plan, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileResourceOwnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileResourceOwnerResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteResourceOwner(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileResourceOwnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileResourceOwner))...)
}

func (r configurationProfileResourceOwnerResource) updateResourceOwnerModel(m *models.ConfigurationProfileResourceOwnerResource, config *api.ConfigurationProfile) {
	m.ID = types.StringValue(config.ID)
	m.Profile = types.StringValue(config.Profile)
	m.Default = tftypes.StringsList(config.Record.ResourceOwnerConfiguration.Default)
	m.OrgDomain = types.StringPointerValue(config.Record.ResourceOwnerConfiguration.OrgDomain)
	m.OrgDomainTag = types.StringPointerValue(config.Record.ResourceOwnerConfiguration.OrgDomainTag)
	m.Tags = tftypes.StringsList(config.Record.ResourceOwnerConfiguration.Tags)
}
