package resources

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &accountGroupItemResource{}
	_ resource.ResourceWithConfigure   = &accountGroupItemResource{}
	_ resource.ResourceWithImportState = &accountGroupItemResource{}
)

func NewAccountGroupItemResource() resource.Resource {
	return &accountGroupItemResource{}
}

type accountGroupItemResource struct {
	api *api.API
}

func (r *accountGroupItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group_item"
}

func (r *accountGroupItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an account within an account group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the account group item.",
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
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *accountGroupItemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *accountGroupItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.AccountGroupItemResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.api.AccountGroupItem.Create(ctx, plan.CloudProvider.ValueString(), plan.AccountKey.ValueString(), plan.GroupUUID.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountGroupItemModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountGroupItemResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupItem, err := r.api.AccountGroupItem.Read(ctx, state.CloudProvider.ValueString(), state.AccountKey.ValueString(), state.GroupUUID.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	}

	if accountGroupItem.ID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updateAccountGroupItemModel(&state, accountGroupItem)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountGroupItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AccountGroupItemResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// There's nothing that can be updated in the state

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AccountGroupItemResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.AccountGroupItem.Delete(ctx, state.ID.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *accountGroupItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: $group_uuid:$cloud_provider:$account_key",
		)
		return
	}

	groupUUID := parts[0]
	cloudProvider := parts[1]
	accountKey := parts[2]

	accountGroupItem, err := r.api.AccountGroupItem.Read(ctx, cloudProvider, accountKey, groupUUID)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), accountGroupItem.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_uuid"), accountGroupItem.GroupUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), accountGroupItem.Provider)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_key"), accountGroupItem.AccountKey)...)
}

func updateAccountGroupItemModel(m *models.AccountGroupItemResource, accountGroupItem api.AccountGroupItem) {
	m.ID = types.StringValue(accountGroupItem.ID)
	m.GroupUUID = types.StringValue(accountGroupItem.GroupUUID)
	m.AccountKey = types.StringValue(accountGroupItem.AccountKey)
	m.CloudProvider = types.StringValue(string(accountGroupItem.Provider))
}
