package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &RepositoryDataSource{}

func NewRepositoryDataSource() datasource.DataSource {
	return &RepositoryDataSource{}
}

// RepositoryDataSource defines the data source implementation.
type RepositoryDataSource struct {
	client *graphql.Client
}

// RepositoryDataSourceModel describes the data source data model.
type RepositoryDataSourceModel struct {
	UUID              types.String   `tfsdk:"uuid"`
	Name              types.String   `tfsdk:"name"`
	URL               types.String   `tfsdk:"url"`
	Description       types.String   `tfsdk:"description"`
	PolicyFileSuffix  []types.String `tfsdk:"policy_file_suffix"`
	PolicyDirectories []types.String `tfsdk:"policy_directories"`
	BranchName        types.String   `tfsdk:"branch_name"`
	AuthUser          types.String   `tfsdk:"auth_user"`
	HasAuthToken      types.Bool     `tfsdk:"has_auth_token"`
	HasSSHPrivateKey  types.Bool     `tfsdk:"has_ssh_private_key"`
	HasSSHPassphrase  types.Bool     `tfsdk:"has_ssh_passphrase"`
	Head              types.String   `tfsdk:"head"`
	LastScanned       types.String   `tfsdk:"last_scanned"`
	VCSProvider       types.String   `tfsdk:"vcs_provider"`
	System            types.Bool     `tfsdk:"system"`
}

func (d *RepositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *RepositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a Stacklet repository.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "The UUID of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the repository.",
				Computed:    true,
			},
			"policy_file_suffix": schema.ListAttribute{
				Description: "The file suffixes used for policy scanning.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"policy_directories": schema.ListAttribute{
				Description: "The directories that are scanned for policies.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"branch_name": schema.StringAttribute{
				Description: "The branch used for scanning policies.",
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
			"has_ssh_private_key": schema.BoolAttribute{
				Description: "Whether the repository has an SSH private key configured.",
				Computed:    true,
			},
			"has_ssh_passphrase": schema.BoolAttribute{
				Description: "Whether the repository has an SSH passphrase configured.",
				Computed:    true,
			},
			"head": schema.StringAttribute{
				Description: "The head commit that was processed.",
				Computed:    true,
			},
			"last_scanned": schema.StringAttribute{
				Description: "The ISO format datetime of when the repo was last scanned.",
				Computed:    true,
			},
			"vcs_provider": schema.StringAttribute{
				Description: "The provider of the repository (e.g., 'github', 'gitlab').",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system repository (not user editable).",
				Computed:    true,
			},
		},
	}
}

func (d *RepositoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

	d.client = client
}

func (d *RepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RepositoryDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare GraphQL query
	var query struct {
		Repository struct {
			UUID              string
			Name              string
			URL               string
			Description       string
			PolicyFileSuffix  []string
			PolicyDirectories []string
			BranchName        string
			AuthUser          string
			HasAuthToken      bool
			HasSshPrivateKey  bool
			HasSshPassphrase  bool
			Head              string
			LastScanned       string
			Provider          string
			System            bool
		} `graphql:"repository(name: $name)"`
	}

	// Prepare variables based on what's provided
	variables := map[string]interface{}{
		"name": data.Name.ValueString(),
	}

	// Execute query
	if err := d.client.Query(ctx, &query, variables); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read repository, got error: %s", err))
		return
	}

	// Map response to model
	data.UUID = types.StringValue(query.Repository.UUID)
	data.Name = types.StringValue(query.Repository.Name)
	data.URL = types.StringValue(query.Repository.URL)
	data.Description = types.StringValue(query.Repository.Description)
	data.BranchName = types.StringValue(query.Repository.BranchName)
	data.AuthUser = types.StringValue(query.Repository.AuthUser)
	data.HasAuthToken = types.BoolValue(query.Repository.HasAuthToken)
	data.HasSSHPrivateKey = types.BoolValue(query.Repository.HasSshPrivateKey)
	data.HasSSHPassphrase = types.BoolValue(query.Repository.HasSshPassphrase)
	data.Head = types.StringValue(query.Repository.Head)
	data.LastScanned = types.StringValue(query.Repository.LastScanned)
	data.VCSProvider = types.StringValue(query.Repository.Provider)
	data.System = types.BoolValue(query.Repository.System)

	// Map policy file suffixes
	data.PolicyFileSuffix = make([]types.String, len(query.Repository.PolicyFileSuffix))
	for i, suffix := range query.Repository.PolicyFileSuffix {
		data.PolicyFileSuffix[i] = types.StringValue(suffix)
	}

	// Map policy directories
	data.PolicyDirectories = make([]types.String, len(query.Repository.PolicyDirectories))
	for i, dir := range query.Repository.PolicyDirectories {
		data.PolicyDirectories[i] = types.StringValue(dir)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
