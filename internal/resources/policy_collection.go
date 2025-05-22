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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
	"github.com/stacklet/terraform-provider-stacklet/schemavalidate"
)

var (
	_ resource.Resource                = &policyCollectionResource{}
	_ resource.ResourceWithConfigure   = &policyCollectionResource{}
	_ resource.ResourceWithImportState = &policyCollectionResource{}
)

func NewPolicyCollectionResource() resource.Resource {
	return &policyCollectionResource{}
}

type policyCollectionResource struct {
	api *api.API
}

func (r *policyCollectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection"
}

func (r *policyCollectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the policy collection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy collection.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the policy collection.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the policy collection (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
				Validators: []validator.String{
					schemavalidate.OneOfCloudProviders(),
				},
			},
			"auto_update": schema.BoolAttribute{
				Description: "Whether the policy collection automatically updates policy versions.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system policy collection.",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dynamic": schema.BoolAttribute{
				Description: "Whether this is a dynamic policy collection.",
				Computed:    true,
			},
			"repository_uuid": schema.StringAttribute{
				Description: "The UUID of the repository the collection is linked to, if dynamic.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *policyCollectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *policyCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PolicyCollectionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionCreateInput{
		Name:        plan.Name.ValueString(),
		Provider:    api.CloudProvider(plan.CloudProvider.ValueString()),
		Description: api.NullableString(plan.Description),
		AutoUpdate:  plan.AutoUpdate.ValueBoolPointer(),
	}
	policyCollection, err := r.api.PolicyCollection.Create(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updatePolicyCollectionModel(&plan, policyCollection)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PolicyCollectionResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyCollection, err := r.api.PolicyCollection.Read(ctx, state.UUID.ValueString(), "")
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	if policyCollection.UUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updatePolicyCollectionModel(&state, policyCollection)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PolicyCollectionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionUpdateInput{
		UUID:        plan.UUID.ValueString(),
		Name:        plan.Name.ValueString(),
		Provider:    api.CloudProvider(plan.CloudProvider.ValueString()),
		Description: plan.Description.ValueStringPointer(),
		AutoUpdate:  plan.AutoUpdate.ValueBoolPointer(),
	}

	policyCollection, err := r.api.PolicyCollection.Update(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updatePolicyCollectionModel(&plan, policyCollection)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PolicyCollectionResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.PolicyCollection.Delete(ctx, state.UUID.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *policyCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), req.ID)...)
}

func updatePolicyCollectionModel(m *models.PolicyCollectionResource, policyCollection api.PolicyCollection) {
	m.ID = types.StringValue(policyCollection.ID)
	m.UUID = types.StringValue(policyCollection.UUID)
	m.Name = types.StringValue(policyCollection.Name)
	m.Description = tftypes.NullableString(policyCollection.Description)
	m.CloudProvider = types.StringValue(string(policyCollection.Provider))
	m.AutoUpdate = types.BoolValue(policyCollection.AutoUpdate)
	m.System = types.BoolValue(policyCollection.System)
	m.Dynamic = types.BoolValue(policyCollection.IsDynamic)
	m.RepositoryUUID = tftypes.NullableString(policyCollection.RepositoryConfig.UUID)
}
