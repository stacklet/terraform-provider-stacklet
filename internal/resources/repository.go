// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &repositoryResource{}
var _ resource.ResourceWithConfigure = &repositoryResource{}
var _ resource.ResourceWithImportState = &repositoryResource{}

func NewRepositoryResource() resource.Resource {
	return &repositoryResource{}
}

// repositoryResource defines the resource implementation.
type repositoryResource struct {
	api *api.API
}

func (r *repositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r *repositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Stacklet repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL node ID.",
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
				Description: "The URL of the remote repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The human-readable name of the repository.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "An optional description of the repository.",
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
				Description: "A URL which triggers scans of the remote repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_user": schema.StringAttribute{
				Description: "The user with access to the remote repository.",
				Optional:    true,
				Computed:    true,
			},
			"has_auth_token": schema.BoolAttribute{
				Description: "Whether auth_token_wo has a value set.",
				Computed:    true,
			},
			"ssh_public_key": schema.StringAttribute{
				Description: "The public key associated with the value set via ssh_private_key_wo.",
				Computed:    true,
			},
			"has_ssh_private_key": schema.BoolAttribute{
				Description: "Whether ssh_private_key_wo has a value set.",
				Computed:    true,
			},
			"has_ssh_passphrase": schema.BoolAttribute{
				Description: "Whether ssh_passphrase_wo has a value set.",
				Computed:    true,
			},
			// After this, write-only secrets and associated trigger attrs.
			"auth_token_wo": schema.StringAttribute{
				Description: "User password/token, or IAM role, to access the remote repository.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"auth_token_wo_version": schema.StringAttribute{
				Description: "Change value to update auth_token_wo.",
				Optional:    true,
			},
			"ssh_private_key_wo": schema.StringAttribute{
				Description: "SSH private key for remote repository authentication.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"ssh_private_key_wo_version": schema.StringAttribute{
				Description: "Change value to update ssh_private_key_wo.",
				Optional:    true,
			},
			"ssh_passphrase_wo": schema.StringAttribute{
				Description: "Passphrase for the SSH private key.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"ssh_passphrase_wo_version": schema.StringAttribute{
				Description: "Change value to update ssh_passphrase_wo.",
				Optional:    true,
			},
		},
	}
}

func (r *repositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read plan and config.
	var plan, config models.RepositoryResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create remote from plan and config.
	auth := api.NewRepositoryConfigAuthInput()
	auth.SetAuthUser(plan.AuthUser.ValueStringPointer())
	auth.SetAuthToken(config.AuthTokenWO.ValueStringPointer())
	auth.SetSSHPrivateKey(config.SSHPrivateKeyWO.ValueStringPointer())
	auth.SetSSHPassphrase(config.SSHPassphraseWO.ValueStringPointer())
	input := api.RepositoryCreateInput{
		Name:        plan.Name.ValueString(),
		URL:         plan.URL.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Auth:        auth,
	}
	repo, err := r.api.Repository.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Update plan from response and save into state.
	r.updateRepositoryModel(&plan, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read known state.
	var state models.RepositoryResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read remote by UUID.
	repo, err := r.api.Repository.Read(ctx, state.UUID.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	// Update state from response.
	r.updateRepositoryModel(&state, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read everything.
	var plan, state, config models.RepositoryResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine which auth fields need to be updated.
	auth := api.NewRepositoryConfigAuthInput()
	auth.SetAuthUser(plan.AuthUser.ValueStringPointer())
	if state.AuthTokenWOVersion != plan.AuthTokenWOVersion {
		state.AuthTokenWOVersion = plan.AuthTokenWOVersion
		auth.SetAuthToken(config.AuthTokenWO.ValueStringPointer())
	}
	if state.SSHPrivateKeyWOVersion != plan.SSHPrivateKeyWOVersion {
		state.SSHPrivateKeyWOVersion = plan.SSHPrivateKeyWOVersion
		auth.SetSSHPrivateKey(config.SSHPrivateKeyWO.ValueStringPointer())
	}
	if state.SSHPassphraseWOVersion != plan.SSHPassphraseWOVersion {
		state.SSHPassphraseWOVersion = plan.SSHPassphraseWOVersion
		auth.SetSSHPassphrase(config.SSHPassphraseWO.ValueStringPointer())
	}

	// Update remote according to the combined grand plan.
	input := api.RepositoryUpdateInput{
		UUID:        plan.UUID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Auth:        auth,
	}
	repo, err := r.api.Repository.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Merge response into state.
	r.updateRepositoryModel(&state, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *repositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	uuid, err := r.api.Repository.FindByURL(ctx, req.ID)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), uuid)...)
}

func (r repositoryResource) updateRepositoryModel(m *models.RepositoryResource, repo *api.Repository) {
	m.ID = types.StringValue(repo.ID)
	m.UUID = types.StringValue(repo.UUID)
	m.URL = types.StringValue(repo.URL)
	m.Name = types.StringValue(repo.Name)
	m.Description = types.StringPointerValue(repo.Description)
	m.System = types.BoolValue(repo.System)
	m.WebhookURL = types.StringValue(repo.WebhookURL)
	m.AuthUser = types.StringPointerValue(repo.Auth.AuthUser)
	m.HasAuthToken = types.BoolValue(repo.Auth.HasAuthToken)
	m.SSHPublicKey = types.StringPointerValue(repo.Auth.SSHPublicKey)
	m.HasSSHPrivateKey = types.BoolValue(repo.Auth.HasSshPrivateKey)
	m.HasSSHPassphrase = types.BoolValue(repo.Auth.HasSshPassphrase)
}
