package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &repositoryDataSource{}

func NewRepositoryDataSource() datasource.DataSource {
	return &repositoryDataSource{}
}

// repositoryDataSource defines the data source implementation.
type repositoryDataSource struct {
	api *api.API
}

func (d *repositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *repositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a Stacklet repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL ID of the repository.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the repository.",
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "The URL of the webhook which triggers repository scans.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the repository.",
				Computed:    true,
			},
			"auth_user": schema.StringAttribute{
				Description: "The user used to access the repository.",
				Computed:    true,
			},
			"has_auth_token": schema.BoolAttribute{
				Description: "Whether the repository has an auth token configured.",
				Computed:    true,
			},
			"ssh_public_key": schema.StringAttribute{
				Description: "If has_ssh_private_key, identifies that SSH private key.",
				Computed:    true,
			},
			"has_ssh_private_key": schema.BoolAttribute{
				Description: "Whether the repository has an SSH private key configured.",
				Computed:    true,
			},
			"has_ssh_passphrase": schema.BoolAttribute{
				Description: "Whether the repository has an SSH passphrase configured.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system repository (not user editable).",
				Computed:    true,
			},
		},
	}
}

func (d *repositoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.api = api.New(client)
}

func (d *repositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.RepositoryDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find by UUID, looking that up by URL if necessary.
	var uuid string
	if data.UUID.IsNull() || data.UUID.IsUnknown() {
		var err error
		uuid, err = d.api.Repository.FindByURL(ctx, data.URL.ValueString())
		if err != nil {
			helpers.AddDiagError(&resp.Diagnostics, err)
			return
		}
	} else {
		uuid = data.UUID.ValueString()
	}
	repo, err := d.api.Repository.Read(ctx, uuid)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Map response to model
	data.ID = types.StringValue(repo.ID)
	data.UUID = types.StringValue(repo.UUID)
	data.URL = types.StringValue(repo.URL)
	data.Name = types.StringValue(repo.Name)
	data.WebhookURL = types.StringValue(repo.WebhookURL)
	data.Description = tftypes.NullableString(repo.Description)
	data.AuthUser = tftypes.NullableString(repo.Auth.AuthUser)
	data.HasAuthToken = types.BoolValue(repo.Auth.HasAuthToken)
	data.SSHPublicKey = tftypes.NullableString(repo.Auth.SSHPublicKey)
	data.HasSSHPrivateKey = types.BoolValue(repo.Auth.HasSshPrivateKey)
	data.HasSSHPassphrase = types.BoolValue(repo.Auth.HasSshPassphrase)
	data.System = types.BoolValue(repo.System)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
