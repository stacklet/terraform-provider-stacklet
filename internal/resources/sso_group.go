package resources

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

var (
	_ resource.Resource                = &ssoGroupResource{}
	_ resource.ResourceWithImportState = &ssoGroupResource{}
)

func NewSSOGroupResource() resource.Resource {
	return &ssoGroupResource{}
}

type ssoGroupResource struct {
	client *graphql.Client
}

type ssoGroupResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Roles             types.List   `tfsdk:"roles"`
	AccountGroupUUIDs types.List   `tfsdk:"account_group_uuids"`
}

type SetSSOGroupConfigsInput struct {
	Groups []ssoGroupConfigInput `json:"groups"`
}

type ssoGroupConfigInput struct {
	Name              string   `json:"name"`
	Roles             []string `json:"roles"`
	AccountGroupUUIDs []string `json:"accountGroupUUIDs"`
}

// Helper function to sort string slices.
func sortStrings(s []string) []string {
	sorted := make([]string, len(s))
	copy(sorted, s)
	sort.Strings(sorted)
	return sorted
}

func (r *ssoGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_group"
}

func (r *ssoGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an SSO group configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for this SSO group configuration.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name that identifies the group in the external SSO provider.",
				Required:    true,
			},
			"roles": schema.ListAttribute{
				Description: "List of Stacklet API roles automatically granted to SSO users in this group.",
				Required:    true,
				ElementType: types.StringType,
			},
			"account_group_uuids": schema.ListAttribute{
				Description: "List of account group UUIDs whose resources are visible to SSO users in this group.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *ssoGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ssoGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ssoGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// First, get all existing SSO groups
	var query struct {
		SSOGroupConfigs []struct {
			Name              string
			Roles             []string
			AccountGroupUUIDs []string `graphql:"accountGroupUUIDs"`
		} `graphql:"ssoGroupConfigs"`
	}

	err := r.client.Query(ctx, &query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read existing SSO groups, got error: %s", err))
		return
	}

	// Convert plan values to native types and sort them
	var roles []string
	resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roles, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	roles = sortStrings(roles)

	var accountGroupUUIDs []string
	resp.Diagnostics.Append(plan.AccountGroupUUIDs.ElementsAs(ctx, &accountGroupUUIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountGroupUUIDs = sortStrings(accountGroupUUIDs)

	// Create a new list of groups with our new group
	groups := make([]ssoGroupConfigInput, 0, len(query.SSOGroupConfigs)+1)
	for _, group := range query.SSOGroupConfigs {
		if group.Name == plan.Name.ValueString() {
			resp.Diagnostics.AddError(
				"Group Already Exists",
				fmt.Sprintf("An SSO group with name %q already exists", plan.Name.ValueString()),
			)
			return
		}
		// Sort the values from existing groups as well
		groups = append(groups, ssoGroupConfigInput{
			Name:              group.Name,
			Roles:             sortStrings(group.Roles),
			AccountGroupUUIDs: sortStrings(group.AccountGroupUUIDs),
		})
	}

	// Add our new group
	groups = append(groups, ssoGroupConfigInput{
		Name:              plan.Name.ValueString(),
		Roles:             roles,
		AccountGroupUUIDs: accountGroupUUIDs,
	})

	// Update all groups
	var mutation struct {
		SetSSOGroups struct {
			Groups []api.SSOGroupConfig `graphql:"groups"`
		} `graphql:"setSSOGroups(input: $input)"`
	}

	variables := map[string]any{
		"input": SetSSOGroupConfigsInput{
			Groups: groups,
		},
	}

	err = r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SSO group, got error: %s", err))
		return
	}

	// Find our group in the response
	var ourGroup *api.SSOGroupConfig
	for _, group := range mutation.SetSSOGroups.Groups {
		if group.Name == plan.Name.ValueString() {
			ourGroup = &group
			break
		}
	}

	if ourGroup == nil {
		resp.Diagnostics.AddError("API Error", "Created group not found in response")
		return
	}

	// Generate a stable ID based on the name
	plan.ID = types.StringValue(fmt.Sprintf("sso-group-%s", ourGroup.Name))
	plan.Name = types.StringValue(ourGroup.Name)

	// Sort the returned values
	roles = sortStrings(ourGroup.Roles)
	accountGroupUUIDs = sortStrings(ourGroup.AccountGroupUUIDs)

	// Create new list values with sorted data
	rolesValue, diags := types.ListValueFrom(ctx, types.StringType, roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Roles = rolesValue

	accountGroupUUIDsValue, diags := types.ListValueFrom(ctx, types.StringType, accountGroupUUIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.AccountGroupUUIDs = accountGroupUUIDsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ssoGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ssoGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		SSOGroupConfigs []api.SSOGroupConfig `graphql:"ssoGroupConfigs"`
	}

	err := r.client.Query(ctx, &query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSO groups, got error: %s", err))
		return
	}

	// Find our group
	var ourGroup *api.SSOGroupConfig
	for _, group := range query.SSOGroupConfigs {
		if group.Name == state.Name.ValueString() {
			ourGroup = &group
			break
		}
	}

	if ourGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Generate a stable ID based on the name
	state.ID = types.StringValue(fmt.Sprintf("sso-group-%s", ourGroup.Name))
	state.Name = types.StringValue(ourGroup.Name)

	// Sort the returned values
	roles := sortStrings(ourGroup.Roles)
	accountGroupUUIDs := sortStrings(ourGroup.AccountGroupUUIDs)

	// Create new list values with sorted data
	rolesValue, diags := types.ListValueFrom(ctx, types.StringType, roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Roles = rolesValue

	accountGroupUUIDsValue, diags := types.ListValueFrom(ctx, types.StringType, accountGroupUUIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.AccountGroupUUIDs = accountGroupUUIDsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ssoGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ssoGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// First, get all existing SSO groups
	var query struct {
		SSOGroupConfigs []struct {
			Name              string
			Roles             []string
			AccountGroupUUIDs []string `graphql:"accountGroupUUIDs"`
		} `graphql:"ssoGroupConfigs"`
	}

	err := r.client.Query(ctx, &query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read existing SSO groups, got error: %s", err))
		return
	}

	// Convert plan values to native types and sort them
	var roles []string
	resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roles, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	roles = sortStrings(roles)

	var accountGroupUUIDs []string
	resp.Diagnostics.Append(plan.AccountGroupUUIDs.ElementsAs(ctx, &accountGroupUUIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountGroupUUIDs = sortStrings(accountGroupUUIDs)

	// Create a new list of groups with our updated group
	groups := make([]ssoGroupConfigInput, 0, len(query.SSOGroupConfigs))
	for _, group := range query.SSOGroupConfigs {
		if group.Name == plan.Name.ValueString() {
			// Skip the old version of our group
			continue
		}
		// Sort the values from existing groups as well
		groups = append(groups, ssoGroupConfigInput{
			Name:              group.Name,
			Roles:             sortStrings(group.Roles),
			AccountGroupUUIDs: sortStrings(group.AccountGroupUUIDs),
		})
	}

	// Add our updated group
	groups = append(groups, ssoGroupConfigInput{
		Name:              plan.Name.ValueString(),
		Roles:             roles,
		AccountGroupUUIDs: accountGroupUUIDs,
	})

	// Update all groups
	var mutation struct {
		SetSSOGroups struct {
			Groups []api.SSOGroupConfig `graphql:"groups"`
		} `graphql:"setSSOGroups(input: $input)"`
	}

	variables := map[string]any{
		"input": SetSSOGroupConfigsInput{
			Groups: groups,
		},
	}

	err = r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update SSO group, got error: %s", err))
		return
	}

	// Find our group in the response
	var ourGroup *api.SSOGroupConfig
	for _, group := range mutation.SetSSOGroups.Groups {
		if group.Name == plan.Name.ValueString() {
			ourGroup = &group
			break
		}
	}

	if ourGroup == nil {
		resp.Diagnostics.AddError("API Error", "Updated group not found in response")
		return
	}

	// Generate a stable ID based on the name
	plan.ID = types.StringValue(fmt.Sprintf("sso-group-%s", ourGroup.Name))
	plan.Name = types.StringValue(ourGroup.Name)

	// Sort the returned values
	roles = sortStrings(ourGroup.Roles)
	accountGroupUUIDs = sortStrings(ourGroup.AccountGroupUUIDs)

	// Create new list values with sorted data
	rolesValue, diags := types.ListValueFrom(ctx, types.StringType, roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Roles = rolesValue

	accountGroupUUIDsValue, diags := types.ListValueFrom(ctx, types.StringType, accountGroupUUIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.AccountGroupUUIDs = accountGroupUUIDsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ssoGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ssoGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// First, get all existing SSO groups
	var query struct {
		SSOGroupConfigs []struct {
			Name              string
			Roles             []string
			AccountGroupUUIDs []string `graphql:"accountGroupUUIDs"`
		} `graphql:"ssoGroupConfigs"`
	}

	err := r.client.Query(ctx, &query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read existing SSO groups, got error: %s", err))
		return
	}

	// Create a new list of groups without our group
	groups := make([]ssoGroupConfigInput, 0, len(query.SSOGroupConfigs)-1)
	for _, group := range query.SSOGroupConfigs {
		if group.Name == state.Name.ValueString() {
			// Skip our group
			continue
		}
		groups = append(groups, ssoGroupConfigInput{
			Name:              group.Name,
			Roles:             group.Roles,
			AccountGroupUUIDs: group.AccountGroupUUIDs,
		})
	}

	// Update all groups
	var mutation struct {
		SetSSOGroups struct {
			Groups []api.SSOGroupConfig `graphql:"groups"`
		} `graphql:"setSSOGroups(input: $input)"`
	}

	variables := map[string]any{
		"input": SetSSOGroupConfigsInput{
			Groups: groups,
		},
	}

	err = r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SSO group, got error: %s", err))
		return
	}
}

func (r *ssoGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is expected to be the name of the SSO group
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
