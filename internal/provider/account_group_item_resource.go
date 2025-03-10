package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ resource.Resource = &accountGroupItemResource{}
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
			},
			"account_key": schema.StringAttribute{
				Description: "The Key of the account to add to the group.",
				Required:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
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

	variables := map[string]interface{}{
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

	// Generate a stable ID for the account group item
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

	variables := map[string]interface{}{
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
	// Account group items don't have any updateable attributes
	var plan accountGroupItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	variables := map[string]interface{}{
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
