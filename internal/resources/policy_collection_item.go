package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ resource.Resource                = &policyCollectionItemResource{}
	_ resource.ResourceWithImportState = &policyCollectionItemResource{}
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
		Description: "Manage a policy within a policy collection. This resource allows to add or remove policies from collections.",
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
							ID     string
							Policy struct {
								UUID string
							}
						}
					}
				}
			}
		} `graphql:"addPolicyCollectionItems(input: $input)"`
	}

	variables := map[string]any{
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

	var mappingID string
	for _, edge := range mutation.AddPolicyCollectionItems.Collection.PolicyMappings.Edges {
		if edge.Node.Policy.UUID == plan.PolicyUUID.ValueString() {
			mappingID = edge.Node.ID
			break
		}
	}
	if mappingID == "" {
		resp.Diagnostics.AddError("Client Error", "Policy collection item was not found in collection after adding")
		return
	}

	plan.ID = types.StringValue(mappingID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyCollectionItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query, err := r.getCollection(ctx, state.CollectionUUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy collection, got error: %s", err))
		return
	}

	// Find the policy in the collection
	var foundPolicy bool
	for _, edge := range query.PolicyCollection.PolicyMappings.Edges {
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

	variables := map[string]any{
		"input": RemovePolicyCollectionMappingsInput{
			IDs: []graphql.ID{graphql.ID(state.ID.ValueString())},
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

func (r *policyCollectionItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: collection_uuid:policy_uuid",
		)
		return
	}

	collectionUUID := parts[0]
	policyUUID := parts[1]

	query, err := r.getCollection(ctx, collectionUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Policy Collection",
			fmt.Sprintf("Could not read policy collection with UUID %s: %s", collectionUUID, err),
		)
		return
	}

	var mappingID string
	for _, edge := range query.PolicyCollection.PolicyMappings.Edges {
		if edge.Node.Policy.UUID == policyUUID {
			mappingID = edge.Node.ID
			break
		}
	}
	if mappingID == "" {
		resp.Diagnostics.AddError("Client Error", "Policy collection item was not found in collection")
		return
	}

	// Set the imported attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_uuid"), collectionUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_uuid"), policyUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), mappingID)...)
}

type collectionQuery struct {
	PolicyCollection struct {
		PolicyMappings struct {
			Edges []struct {
				Node struct {
					ID     string
					Policy struct {
						UUID string
					}
				}
			}
		}
	} `graphql:"policyCollection(uuid: $uuid)"`
}

func (r *policyCollectionItemResource) getCollection(ctx context.Context, uuid string) (collectionQuery, error) {
	var query collectionQuery
	variables := map[string]any{
		"uuid": graphql.String(uuid),
	}
	err := r.client.Query(ctx, &query, variables)
	return query, err
}
