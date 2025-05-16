package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ resource.Resource                = &policyCollectionItemResource{}
	_ resource.ResourceWithConfigure   = &policyCollectionItemResource{}
	_ resource.ResourceWithImportState = &policyCollectionItemResource{}
)

func NewPolicyCollectionItemResource() resource.Resource {
	return &policyCollectionItemResource{}
}

type policyCollectionItemResource struct {
	api *api.API
}

func (r *policyCollectionItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection_item"
}

func (r *policyCollectionItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a policy within a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the policy collection item.",
				Computed:    true,
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
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *policyCollectionItemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyCollectionItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PolicyCollectionItemResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionItemCreateInput{
		CollectionUUID: plan.CollectionUUID.ValueString(),
		PolicyUUID:     plan.PolicyUUID.ValueString(),
		PolicyVersion:  int(plan.PolicyVersion.ValueInt32()),
	}
	policyCollectionItem, err := r.api.PolicyCollectionItem.Create(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updatePolicyCollectionItemModel(&plan, policyCollectionItem)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PolicyCollectionItemResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyCollectionItem, err := r.api.PolicyCollectionItem.Read(ctx, state.CollectionUUID.ValueString(), state.PolicyUUID.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	}

	if policyCollectionItem.ID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updatePolicyCollectionItemModel(&state, policyCollectionItem)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyCollectionItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PolicyCollectionItemResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PolicyCollectionItems don't have any updatable attributes

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PolicyCollectionItemResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.PolicyCollectionItem.Delete(ctx, state.ID.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *policyCollectionItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := helpers.SplitImportID(req.ID, []string{"collection_uuid", "policy_uuid"})
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_uuid"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_uuid"), parts[1])...)
}

func updatePolicyCollectionItemModel(m *models.PolicyCollectionItemResource, policyCollectionItem api.PolicyCollectionItem) {
	m.ID = types.StringValue(policyCollectionItem.ID)
	m.CollectionUUID = types.StringValue(policyCollectionItem.Collection.UUID)
	m.PolicyUUID = types.StringValue(policyCollectionItem.Policy.UUID)
	m.PolicyVersion = types.Int32Value(int32(policyCollectionItem.Policy.Version))
}
