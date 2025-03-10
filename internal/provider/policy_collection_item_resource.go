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
	_ resource.Resource = &policyCollectionItemResource{}
)

func NewPolicyCollectionItemResource() resource.Resource {
	return &policyCollectionItemResource{}
}

type policyCollectionItemResource struct {
	client *graphql.Client
}

type policyCollectionItemResourceModel struct {
	ID             types.String `tfsdk:"id"`
	CollectionUUID types.String `tfsdk:"collection_uuid"`
	PolicyUUID     types.String `tfsdk:"policy_uuid"`
}

func (r *policyCollectionItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection_item"
}

func (r *policyCollectionItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a policy within a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the policy collection item.",
				Computed:    true,
			},
			"collection_uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection.",
				Required:    true,
			},
			"policy_uuid": schema.StringAttribute{
				Description: "The UUID of the policy to add to the collection.",
				Required:    true,
			},
		},
	}
}

func (r *policyCollectionItemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyCollectionItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyCollectionItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		AddPolicyCollectionItems struct {
			Collection struct {
				UUID           string
				PolicyMappings struct {
					Edges []struct {
						Node struct {
							ID string
						}
					}
				}
			}
		} `graphql:"addPolicyCollectionItems(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": PolicyCollectionItemsInput{
			UUID: plan.CollectionUUID.ValueString(),
			Items: []PolicyCollectionElement{
				{
					PolicyUUID: plan.PolicyUUID.ValueString(),
				},
			},
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add policy to collection, got error: %s", err))
		return
	}

	// Find the added policy in the response
	var addedPolicy bool
	for _, edge := range mutation.AddPolicyCollectionItems.Collection.PolicyMappings.Edges {
		if edge.Node.ID != "" {
			addedPolicy = true
			break
		}
	}

	if !addedPolicy {
		resp.Diagnostics.AddError("Client Error", "Policy was not found in collection after adding")
		return
	}

	// Generate a stable ID for the policy collection item
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.CollectionUUID.ValueString(), plan.PolicyUUID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyCollectionItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		PolicyCollection struct {
			Policies struct {
				Edges []struct {
					Node struct {
						Policy struct {
							UUID string
						}
					}
				}
			}
		} `graphql:"policyCollection(uuid: $uuid)"`
	}

	variables := map[string]interface{}{
		"uuid": graphql.String(state.CollectionUUID.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy collection, got error: %s", err))
		return
	}

	// Find the policy in the collection
	var foundPolicy bool
	for _, edge := range query.PolicyCollection.Policies.Edges {
		if edge.Node.Policy.UUID == state.PolicyUUID.ValueString() {
			foundPolicy = true
			break
		}
	}

	if !foundPolicy {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyCollectionItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Policy collection items don't have any updateable attributes
	var plan policyCollectionItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyCollectionItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		RemovePolicyCollectionMappings struct {
			Removed []struct {
				ID string
			}
		} `graphql:"removePolicyCollectionMappings(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": RemovePolicyCollectionMappingsInput{
			IDs: []graphql.ID{wrapNodeID([]string{"policy-collection-mapping", state.CollectionUUID.ValueString(), state.PolicyUUID.ValueString()})},
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove policy from collection, got error: %s", err))
		return
	}
}

// Input types for GraphQL mutations
type PolicyCollectionElement struct {
	PolicyUUID string `json:"policyUUID"`
}

type PolicyCollectionItemsInput struct {
	UUID  string                    `json:"uuid"`
	Items []PolicyCollectionElement `json:"items"`
}

type RemovePolicyCollectionMappingsInput struct {
	IDs []graphql.ID `json:"ids"`
}
