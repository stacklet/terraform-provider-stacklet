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
var _ resource.ResourceWithConfigure = &RepositoryResource{}
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
	// Read plan and config.
	var plan, config models.RepositoryResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create remote from plan and config.
	auth := api.NewRepositoryConfigAuthInput()
	auth.SetAuthUser(api.NullableString(plan.AuthUser))
	auth.SetAuthToken(api.NullableString(config.AuthTokenWO))
	auth.SetSSHPrivateKey(api.NullableString(config.SSHPrivateKeyWO))
	auth.SetSSHPassphrase(api.NullableString(config.SSHPassphraseWO))
	input := api.RepositoryCreateInput{
		Name:        plan.Name.ValueString(),
		URL:         plan.URL.ValueString(),
		Description: api.NullableString(plan.Description),
		Auth:        auth,
	}
	repo, err := r.api.Repository.Create(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Update plan from response and save into state.
	updateRepositoryModel(&plan, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read known state.
	var state models.RepositoryResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read remote by UUID.
	repo, err := r.api.Repository.Read(ctx, state.UUID.ValueString())
	if err != nil {
		if _, ok := err.(api.NotFound); ok {
			resp.State.RemoveResource(ctx)
		} else {
			helpers.AddDiagError(&resp.Diagnostics, err)
		}
		return
	}

	// Update state from response.
	updateRepositoryModel(&state, repo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	auth.SetAuthUser(api.NullableString(plan.AuthUser))
	if state.AuthTokenWOVersion != plan.AuthTokenWOVersion {
		state.AuthTokenWOVersion = plan.AuthTokenWOVersion
		auth.SetAuthToken(api.NullableString(config.AuthTokenWO))
	}
	if state.SSHPrivateKeyWOVersion != plan.SSHPrivateKeyWOVersion {
		state.SSHPrivateKeyWOVersion = plan.SSHPrivateKeyWOVersion
		auth.SetSSHPrivateKey(api.NullableString(config.SSHPrivateKeyWO))
	}
	if state.SSHPassphraseWOVersion != plan.SSHPassphraseWOVersion {
		state.SSHPassphraseWOVersion = plan.SSHPassphraseWOVersion
		auth.SetSSHPassphrase(api.NullableString(config.SSHPassphraseWO))
	}

	// Update remote according to the combined grand plan.
	input := api.RepositoryUpdateInput{
		UUID:        plan.UUID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: api.NullableString(plan.Description),
		Auth:        auth,
	}
	repo, err := r.api.Repository.Update(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Merge response into state.
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
	uuid, err := r.api.Repository.FindByURL(ctx, req.ID)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), uuid)...)
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
	m.SSHPublicKey = tftypes.NullableString(repo.Auth.SSHPublicKey)
	m.HasSSHPrivateKey = types.BoolValue(repo.Auth.HasSshPrivateKey)
	m.HasSSHPassphrase = types.BoolValue(repo.Auth.HasSshPassphrase)
}
