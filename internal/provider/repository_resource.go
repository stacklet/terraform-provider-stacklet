package provider

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
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RepositoryResource{}
var _ resource.ResourceWithImportState = &RepositoryResource{}

func NewRepositoryResource() resource.Resource {
	return &RepositoryResource{}
}

// RepositoryResource defines the resource implementation.
type RepositoryResource struct {
	client *graphql.Client
}

// RepositoryResourceModel describes the resource data model.
type RepositoryResourceModel struct {
	UUID              types.String   `tfsdk:"uuid"`
	Name              types.String   `tfsdk:"name"`
	URL               types.String   `tfsdk:"url"`
	Description       types.String   `tfsdk:"description"`
	PolicyFileSuffix  []types.String `tfsdk:"policy_file_suffix"`
	PolicyDirectories []types.String `tfsdk:"policy_directories"`
	BranchName        types.String   `tfsdk:"branch_name"`
	AuthUser          types.String   `tfsdk:"auth_user"`
	AuthToken         types.String   `tfsdk:"auth_token"`
	SSHPrivateKey     types.String   `tfsdk:"ssh_private_key"`
	SSHPassphrase     types.String   `tfsdk:"ssh_passphrase"`
	DeepImport        types.Bool     `tfsdk:"deep_import"`
}

func (r *RepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r *RepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Stacklet repository.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "The UUID of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the repository.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description of the repository.",
				Optional:    true,
			},
			"policy_file_suffix": schema.ListAttribute{
				Description: "Override the default suffix options ['.yaml', '.yml']. This could allow specifying ['.json'] to process other files.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"policy_directories": schema.ListAttribute{
				Description: "If set, only directories that match the list will be scanned for policies.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"branch_name": schema.StringAttribute{
				Description: "If set, use the specified branch name when scanning for policies rather than the repository default.",
				Optional:    true,
			},
			"auth_user": schema.StringAttribute{
				Description: "The user to use to access the repository if it is private.",
				Optional:    true,
			},
			"auth_token": schema.StringAttribute{
				Description: "The token for the user to use to access the repository if it is private.",
				Optional:    true,
				Sensitive:   true,
			},
			"ssh_private_key": schema.StringAttribute{
				Description: "SSH private key for repository authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"ssh_passphrase": schema.StringAttribute{
				Description: "Passphrase for the SSH private key.",
				Optional:    true,
				Sensitive:   true,
			},
			"deep_import": schema.BoolAttribute{
				Description: "If true, scan repository from the beginning. If false, only scan the tip.",
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

	r.client = client
}

func (r *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare GraphQL mutation
	var mutation struct {
		AddRepository struct {
			Repository struct {
				UUID string
			}
		} `graphql:"addRepository(input: $input)"`
	}

	// Convert policy file suffixes
	var policySuffixes []string
	for _, suffix := range data.PolicyFileSuffix {
		policySuffixes = append(policySuffixes, suffix.ValueString())
	}

	// Convert policy directories
	var policyDirs []string
	for _, dir := range data.PolicyDirectories {
		policyDirs = append(policyDirs, dir.ValueString())
	}

	// Prepare variables
	variables := map[string]interface{}{
		"input": RepositoryInput{
			Name:              data.Name.ValueString(),
			URL:               data.URL.ValueString(),
			Description:       toString(data.Description),
			PolicyFileSuffix:  policySuffixes,
			PolicyDirectories: policyDirs,
			BranchName:        toString(data.BranchName),
			AuthUser:          toString(data.AuthUser),
			AuthToken:         toString(data.AuthToken),
			SSHPrivateKey:     toString(data.SSHPrivateKey),
			SSHPassphrase:     toString(data.SSHPassphrase),
			DeepImport:        toBoolPtr(data.DeepImport),
		},
	}

	// Execute mutation
	if err := r.client.Mutate(ctx, &mutation, variables); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create repository, got error: %s", err))
		return
	}

	// Save UUID from response
	data.UUID = types.StringValue(mutation.AddRepository.Repository.UUID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
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
		} `graphql:"repository(url: $url)"`
	}

	// Execute query
	variables := map[string]interface{}{
		"url": data.URL.ValueString(),
	}

	if err := r.client.Query(ctx, &query, variables); err != nil {
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare GraphQL mutation
	var mutation struct {
		UpdateRepository struct {
			Repository struct {
				UUID string
			}
		} `graphql:"updateRepository(input: $input)"`
	}

	// Convert policy file suffixes
	var policySuffixes []string
	for _, suffix := range data.PolicyFileSuffix {
		policySuffixes = append(policySuffixes, suffix.ValueString())
	}

	// Convert policy directories
	var policyDirs []string
	for _, dir := range data.PolicyDirectories {
		policyDirs = append(policyDirs, dir.ValueString())
	}

	// Prepare variables
	variables := map[string]interface{}{
		"input": UpdateRepositoryInput{
			URL:               data.URL.ValueString(),
			Name:              toString(data.Name),
			Description:       toString(data.Description),
			PolicyFileSuffix:  policySuffixes,
			PolicyDirectories: policyDirs,
			BranchName:        toString(data.BranchName),
			AuthUser:          toString(data.AuthUser),
			AuthToken:         toString(data.AuthToken),
			SSHPrivateKey:     toString(data.SSHPrivateKey),
			SSHPassphrase:     toString(data.SSHPassphrase),
		},
	}

	// Execute mutation
	if err := r.client.Mutate(ctx, &mutation, variables); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update repository, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare GraphQL mutation
	var mutation struct {
		RemoveRepository struct {
			Repository struct {
				UUID string
			}
		} `graphql:"removeRepository(url: $url)"`
	}

	// Execute mutation
	variables := map[string]interface{}{
		"url": data.URL.ValueString(),
	}

	if err := r.client.Mutate(ctx, &mutation, variables); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete repository, got error: %s", err))
		return
	}
}

func (r *RepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("url"), req, resp)
}

// Helper functions for handling optional values
func toString(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	value := v.ValueString()
	return &value
}

func toBoolPtr(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	value := v.ValueBool()
	return &value
}

// GraphQL input types
type RepositoryInput struct {
	Name              string   `json:"name"`
	URL               string   `json:"url"`
	Description       *string  `json:"description,omitempty"`
	PolicyFileSuffix  []string `json:"policyFileSuffix,omitempty"`
	PolicyDirectories []string `json:"policyDirectories,omitempty"`
	BranchName        *string  `json:"branchName,omitempty"`
	AuthUser          *string  `json:"authUser,omitempty"`
	AuthToken         *string  `json:"authToken,omitempty"`
	SSHPrivateKey     *string  `json:"sshPrivateKey,omitempty"`
	SSHPassphrase     *string  `json:"sshPassphrase,omitempty"`
	DeepImport        *bool    `json:"deepImport,omitempty"`
}

type UpdateRepositoryInput struct {
	URL               string   `json:"url"`
	Name              *string  `json:"name,omitempty"`
	Description       *string  `json:"description,omitempty"`
	PolicyFileSuffix  []string `json:"policyFileSuffix,omitempty"`
	PolicyDirectories []string `json:"policyDirectories,omitempty"`
	BranchName        *string  `json:"branchName,omitempty"`
	AuthUser          *string  `json:"authUser,omitempty"`
	AuthToken         *string  `json:"authToken,omitempty"`
	SSHPrivateKey     *string  `json:"sshPrivateKey,omitempty"`
	SSHPassphrase     *string  `json:"sshPassphrase,omitempty"`
}
