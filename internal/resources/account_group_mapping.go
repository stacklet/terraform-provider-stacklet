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
	_ resource.Resource                = &accountGroupMappingResource{}
	_ resource.ResourceWithConfigure   = &accountGroupMappingResource{}
	_ resource.ResourceWithImportState = &accountGroupMappingResource{}
)

func NewAccountGroupMappingResource() resource.Resource {
	return &accountGroupMappingResource{}
}

type accountGroupMappingResource struct {
	api *api.API
}

func (r *accountGroupMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group_mapping"
}

func (r *accountGroupMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an account within an account group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the account group mapping.",
				Computed:    true,
			},
			"group_uuid": schema.StringAttribute{
				Description: "The UUID of the account group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_key": schema.StringAttribute{
				Description: "The Key of the account to add to the group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *accountGroupMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *accountGroupMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.AccountGroupMappingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.api.AccountGroupMapping.Create(ctx, plan.AccountKey.ValueString(), plan.GroupUUID.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountGroupMappingModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountGroupMappingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupMapping, err := r.api.AccountGroupMapping.Read(ctx, state.AccountKey.ValueString(), state.GroupUUID.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	}

	if accountGroupMapping.ID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updateAccountGroupMappingModel(&state, accountGroupMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountGroupMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AccountGroupMappingResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// There's nothing that can be updated in the state, as no currently exposed field can be changed.
	// This might change if we end up exposing fields like `regions`.

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AccountGroupMappingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.AccountGroupMapping.Delete(ctx, state.ID.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *accountGroupMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := helpers.SplitImportID(req.ID, []string{"group_uuid", "account_key"})
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_uuid"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_key"), parts[1])...)
}

func updateAccountGroupMappingModel(m *models.AccountGroupMappingResource, accountGroupMapping api.AccountGroupMapping) {
	m.ID = types.StringValue(accountGroupMapping.ID)
	m.GroupUUID = types.StringValue(accountGroupMapping.GroupUUID)
	m.AccountKey = types.StringValue(accountGroupMapping.AccountKey)
}
