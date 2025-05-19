package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RepositoryResource{}
var _ resource.ResourceWithImportState = &RepositoryResource{}

func NewRepositoryResource() resource.Resource {
	return &RepositoryResource{}
}

// RepositoryResource defines the resource implementation.
type RepositoryResource struct {
	api *api.API
}

func (r *RepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r *RepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Stacklet repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The URL of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the repository.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the repository.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				Description: "System repositories cannot be changed.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"webhook_url": schema.StringAttribute{
				Description: "A description of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_user": schema.StringAttribute{
				Description: "The user to use to access the repository if it is private.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"has_auth_token": schema.BoolAttribute{
				Description: "Whether auth_token_wo has a value set.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_ssh_private_key": schema.BoolAttribute{
				Description: "Whether ssh_private_key_wo has a value set.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_ssh_passphrase": schema.BoolAttribute{
				Description: "Whether ssh_passphrase_wo has a value set.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			// After this, write-only secrets and associated trigger attrs.
			"auth_token_wo": schema.StringAttribute{
				Description: "The token for the user to use to access the repository if it is private.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"auth_token_wo_version": schema.Int32Attribute{
				Description: "Change value to update auth_token_wo.",
				Optional:    true,
			},
			"ssh_private_key_wo": schema.StringAttribute{
				Description: "SSH private key for repository authentication.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"ssh_private_key_wo_version": schema.Int32Attribute{
				Description: "Change value to update ssh_private_key_wo.",
				Optional:    true,
			},
			"ssh_passphrase_wo": schema.StringAttribute{
				Description: "Passphrase for the SSH private key.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"ssh_passphrase_wo_version": schema.Int32Attribute{
				Description: "Change value to update ssh_passphrase_wo.",
				Optional:    true,
			},
		},
	}
}

func (r *RepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read plan
	var plan models.RepositoryResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create remote from plan
	input := api.RepositoryCreateInput{
		Name:        plan.Name.ValueString(),
		URL:         plan.URL.ValueString(),
		Description: api.NullableString(plan.Description),
	}
	repo, err := r.api.Repository.Create(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Save response into state
	updateRepositoryModel(&plan, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read state for UUID
	var data models.RepositoryResource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read remote by UUID
	repo, err := r.api.Repository.Read(ctx, data.UUID.ValueString())
	if err != nil {
		if _, ok := err.(api.NotFound); ok {
			resp.State.RemoveResource(ctx)
		} else {
			helpers.AddDiagError(&resp.Diagnostics, err)
		}
		return
	}

	// Save response into state
	updateRepositoryModel(&data, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read plan
	var plan models.RepositoryResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update remote according to plan
	input := api.RepositoryUpdateInput{
		UUID:        plan.UUID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: api.NullableString(plan.Description),
	}
	repo, err := r.api.Repository.Update(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Save response into state
	var state models.RepositoryResource
	updateRepositoryModel(&state, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read state
	var data models.RepositoryResource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete remote
	input := api.RepositoryDeleteInput{
		UUID: data.UUID.ValueString(),
		// Note NOT cascading; we could add a force_delete attr if necessary, but the default
		// behaviour should *NOT* be to implicitly tear down resources not under management.
	}
	if err := r.api.Repository.Delete(ctx, input); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *RepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	repo, err := r.api.Repository.ReadURL(ctx, req.ID)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), repo.UUID)...)
}

func updateRepositoryModel(m *models.RepositoryResource, repo api.Repository) {
	m.ID = types.StringValue(repo.ID)
	m.UUID = types.StringValue(repo.UUID)
	m.URL = types.StringValue(repo.URL)
	m.Name = types.StringValue(repo.Name)
	m.Description = tftypes.NullableString(repo.Description)
	m.System = types.BoolValue(repo.System)
	m.WebhookURL = types.StringValue(repo.WebhookURL)
	m.AuthUser = tftypes.NullableString(repo.Auth.AuthUser)
	m.HasAuthToken = types.BoolValue(repo.Auth.HasAuthToken)
	m.HasSSHPrivateKey = types.BoolValue(repo.Auth.HasSshPrivateKey)
	m.HasSSHPassphrase = types.BoolValue(repo.Auth.HasSshPassphrase)
}
