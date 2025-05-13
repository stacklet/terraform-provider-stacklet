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
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &accountResource{}
	_ resource.ResourceWithConfigure   = &accountResource{}
	_ resource.ResourceWithImportState = &accountResource{}
)

func NewAccountResource() resource.Resource {
	return &accountResource{}
}

type accountResource struct {
	api *api.API
}

func (r *accountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (r *accountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a Stacklet account with a specific cloud provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The cloud specific identifier for the account (e.g., AWS account ID, GCP project ID, Azure subscription UUID).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name for the account.",
				Required:    true,
			},
			"short_name": schema.StringAttribute{
				Description: "The short name for the account.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "More detailed information about the account.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				Description: "The path used to group accounts in a hierarchy.",
				Optional:    true,
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email contact address for the account.",
				Optional:    true,
			},
			"security_context": schema.StringAttribute{
				Description: "The security context for the account.",
				Optional:    true,
				Sensitive:   true,
			},
			"variables": schema.StringAttribute{
				Description: "JSON encoded dict of values used for policy templating.",
				Optional:    true,
			},
		},
	}
}

func (r *accountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.AccountResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := api.NewCloudProvider(plan.CloudProvider.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	input := api.AccountCreateInput{
		Name:            plan.Name.ValueString(),
		Key:             plan.Key.ValueString(),
		Provider:        provider,
		ShortName:       api.NullableString(plan.ShortName),
		Description:     api.NullableString(plan.Description),
		Email:           api.NullableString(plan.Email),
		SecurityContext: api.NullableString(plan.SecurityContext),
		Variables:       api.NullableString(plan.Variables),
	}
	account, err := r.api.Account.Create(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountModel(&plan, account)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := r.api.Account.Read(ctx, state.CloudProvider.ValueString(), state.Key.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	if account.Key == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	updateAccountModel(&state, account)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AccountResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := api.NewCloudProvider(plan.CloudProvider.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	input := api.AccountUpdateInput{
		Key:             plan.Key.ValueString(),
		Provider:        provider,
		Name:            api.NullableString(plan.Name),
		ShortName:       api.NullableString(plan.ShortName),
		Description:     api.NullableString(plan.Description),
		Email:           api.NullableString(plan.Email),
		SecurityContext: api.NullableString(plan.SecurityContext),
		Variables:       api.NullableString(plan.Variables),
	}

	account, err := r.api.Account.Update(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	updateAccountModel(&plan, account)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AccountResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.Account.Delete(ctx, state.CloudProvider.ValueString(), state.Key.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Extract the provider and key from the import ID (format: provider:key)
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: provider:key (e.g., aws:123456789012)",
		)
		return
	}

	provider := idParts[0]
	key := idParts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), provider)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
}

func updateAccountModel(m *models.AccountResource, account api.Account) {
	m.ID = types.StringValue(account.ID)
	m.Key = types.StringValue(account.Key)
	m.Name = types.StringValue(account.Name)
	m.ShortName = tftypes.NullableString(account.ShortName)
	m.CloudProvider = types.StringValue(string(account.Provider))
	m.Description = tftypes.NullableString(account.Description)
	m.Path = tftypes.NullableString(account.Path)
	m.Email = tftypes.NullableString(account.Email)
	m.Variables = tftypes.NullableString(account.Variables)
}
