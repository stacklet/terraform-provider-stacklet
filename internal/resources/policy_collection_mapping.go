package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &policyCollectionMappingResource{}
	_ resource.ResourceWithConfigure   = &policyCollectionMappingResource{}
	_ resource.ResourceWithImportState = &policyCollectionMappingResource{}
)

func NewPolicyCollectionMappingResource() resource.Resource {
	return &policyCollectionMappingResource{}
}

type policyCollectionMappingResource struct {
	api *api.API
}

func (r *policyCollectionMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection_mapping"
}

func (r *policyCollectionMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a policy within a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the policy collection mapping.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"collection_uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_uuid": schema.StringAttribute{
				Description: "The UUID of the policy to add to the collection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_version": schema.Int32Attribute{
				Description: "The version of the policy to add to the collection.",
				Required:    true,
			},
		},
	}
}

func (r *policyCollectionMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.api = api.New(client)
}

func (r *policyCollectionMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PolicyCollectionMappingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionMappingInput{
		CollectionUUID: plan.CollectionUUID.ValueString(),
		PolicyUUID:     plan.PolicyUUID.ValueString(),
		PolicyVersion:  int(plan.PolicyVersion.ValueInt32()),
	}
	// Note that given this is an upsert operation, if the mapping already exists it will be updated.
	policyCollectionMapping, err := r.api.PolicyCollectionMapping.Upsert(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updatePolicyCollectionMappingModel(&plan, policyCollectionMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PolicyCollectionMappingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyCollectionMapping, err := r.api.PolicyCollectionMapping.Read(ctx, state.CollectionUUID.ValueString(), state.PolicyUUID.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	}

	if policyCollectionMapping.ID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updatePolicyCollectionMappingModel(&state, policyCollectionMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyCollectionMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PolicyCollectionMappingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionMappingInput{
		CollectionUUID: plan.CollectionUUID.ValueString(),
		PolicyUUID:     plan.PolicyUUID.ValueString(),
		PolicyVersion:  int(plan.PolicyVersion.ValueInt32()),
	}
	policyCollectionMapping, err := r.api.PolicyCollectionMapping.Upsert(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updatePolicyCollectionMappingModel(&plan, policyCollectionMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PolicyCollectionMappingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.PolicyCollectionMapping.Delete(ctx, state.ID.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *policyCollectionMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := helpers.SplitImportID(req.ID, []string{"collection_uuid", "policy_uuid"})
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_uuid"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_uuid"), parts[1])...)
}

func updatePolicyCollectionMappingModel(m *models.PolicyCollectionMappingResource, policyCollectionMapping api.PolicyCollectionMapping) {
	m.ID = types.StringValue(policyCollectionMapping.ID)
	m.CollectionUUID = types.StringValue(policyCollectionMapping.Collection.UUID)
	m.PolicyUUID = types.StringValue(policyCollectionMapping.Policy.UUID)
	m.PolicyVersion = types.Int32Value(int32(policyCollectionMapping.Policy.Version))
}
