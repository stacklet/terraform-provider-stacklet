package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ resource.Resource                = &accountGroupItemResource{}
	_ resource.ResourceWithImportState = &accountGroupItemResource{}
)

func NewAccountGroupItemResource() resource.Resource {
	return &accountGroupItemResource{}
}

type accountGroupItemResource struct {
	client *graphql.Client
}

type accountGroupItemResourceModel struct {
	ID            types.String `tfsdk:"id"`
	GroupUUID     types.String `tfsdk:"group_uuid"`
	AccountKey    types.String `tfsdk:"account_key"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
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

	r.client = client
}

// generateStableID creates a stable hash from the resource's attributes
func generateStableID(groupUUID, cloudProvider, accountKey string) string {
	// Create a deterministic string to hash
	input := fmt.Sprintf("%s:%s:%s", groupUUID, cloudProvider, accountKey)
	// Create a new SHA256 hash
	hash := sha256.New()
	hash.Write([]byte(input))
	// Get the first 8 characters of the hex-encoded hash
	return hex.EncodeToString(hash.Sum(nil))[:8]
}

func (r *accountGroupItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accountGroupItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		AddAccountGroupItems struct {
			Group struct {
				ID       string
				Accounts struct {
					Edges []struct {
						Node struct {
							Account struct {
								Key string
							}
						}
					}
				}
			}
		} `graphql:"addAccountGroupItems(input: $input)"`
	}

	variables := map[string]any{
		"input": AccountGroupItemsInput{
			UUID: plan.GroupUUID.ValueString(),
			Items: []AccountGroupElement{
				{
					Key:      plan.AccountKey.ValueString(),
					Provider: plan.CloudProvider.ValueString(),
				},
			},
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add account to group, got error: %s", err))
		return
	}

	// Find the added account in the response
	var addedAccount bool
	for _, edge := range mutation.AddAccountGroupItems.Group.Accounts.Edges {
		if edge.Node.Account.Key == plan.AccountKey.ValueString() {
			addedAccount = true
			break
		}
	}

	if !addedAccount {
		resp.Diagnostics.AddError("Client Error", "Account was not found in group after adding")
		return
	}

	// Generate a simple ID for the account group item
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.GroupUUID.ValueString(), plan.AccountKey.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountGroupItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		AccountGroup struct {
			Accounts struct {
				Edges []struct {
					Node struct {
						Account struct {
							Key string
						}
					}
				}
			}
		} `graphql:"accountGroup(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(state.GroupUUID.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account group, got error: %s", err))
		return
	}

	// Find the account in the group
	var foundAccount bool
	for _, edge := range query.AccountGroup.Accounts.Edges {
		if edge.Node.Account.Key == state.AccountKey.ValueString() {
			foundAccount = true
			break
		}
	}

	if !foundAccount {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountGroupItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accountGroupItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a stable ID for the account group item
	plan.ID = types.StringValue(generateStableID(
		plan.GroupUUID.ValueString(),
		plan.CloudProvider.ValueString(),
		plan.AccountKey.ValueString(),
	))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accountGroupItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		RemoveAccountGroupMappings struct {
			Removed []struct {
				ID graphql.ID
			}
		} `graphql:"removeAccountGroupMappings(input: $input)"`
	}

	// Construct the node ID using the proper format
	nodeID := wrapNodeID([]string{
		"account-group-mapping",
		state.GroupUUID.ValueString(),
		state.AccountKey.ValueString(),
	})

	variables := map[string]any{
		"input": RemoveAccountGroupMappingsInput{
			IDs: []graphql.ID{nodeID},
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove account from group, got error: %s", err))
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

	// Query to verify the account group exists and contains the account
	var query struct {
		AccountGroup struct {
			Accounts struct {
				Edges []struct {
					Node struct {
						Account struct {
							Key string
						}
					}
				}
			}
		} `graphql:"accountGroup(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(groupUUID),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Account Group",
			fmt.Sprintf("Could not read account group with UUID %s: %s", groupUUID, err),
		)
		return
	}

	// Verify the account exists in the group
	var accountFound bool
	for _, edge := range query.AccountGroup.Accounts.Edges {
		if edge.Node.Account.Key == accountKey {
			accountFound = true
			break
		}
	}

	if !accountFound {
		resp.Diagnostics.AddError(
			"Account Not Found",
			fmt.Sprintf("Account with key %s was not found in account group %s", accountKey, groupUUID),
		)
		return
	}

	// Set the imported attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_uuid"), groupUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), cloudProvider)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_key"), accountKey)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("%s:%s", groupUUID, accountKey))...)
}

type AccountGroupElement struct {
	Key      string   `json:"key"`
	Provider string   `json:"provider"`
	Regions  []string `json:"regions,omitempty"`
}

type AccountGroupItemsInput struct {
	UUID  string                `json:"uuid"`
	Items []AccountGroupElement `json:"items"`
}

type RemoveAccountGroupMappingsInput struct {
	IDs []graphql.ID `json:"ids"`
}
