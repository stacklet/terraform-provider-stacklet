package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ resource.Resource                = &policyCollectionResource{}
	_ resource.ResourceWithImportState = &policyCollectionResource{}
)

func NewPolicyCollectionResource() resource.Resource {
	return &policyCollectionResource{}
}

type policyCollectionResource struct {
	client *graphql.Client
}

type policyCollectionResourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	AutoUpdate    types.Bool   `tfsdk:"auto_update"`
	System        types.Bool   `tfsdk:"system"`
	Repository    types.String `tfsdk:"repository"`
}

func (r *policyCollectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection"
}

func (r *policyCollectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the policy collection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy collection.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the policy collection.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the policy collection (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
			},
			"auto_update": schema.BoolAttribute{
				Description: "Whether the policy collection automatically updates policy versions.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system policy collection.",
				Computed:    true,
			},
			"repository": schema.StringAttribute{
				Description: "The repository URL if this collection was created from a repo control file.",
				Computed:    true,
			},
		},
	}
}

func (r *policyCollectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyCollectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		AddPolicyCollection struct {
			Collection struct {
				ID          string
				UUID        string
				Name        string
				Description string
				Provider    string
				AutoUpdate  bool
				System      bool
				Repository  string
			}
		} `graphql:"addPolicyCollection(input: $input)"`
	}

	variables := map[string]any{
		"input": AddPolicyCollectionInput{
			Name:        plan.Name.ValueString(),
			Provider:    plan.CloudProvider.ValueString(),
			Description: graphql.String(plan.Description.ValueString()),
			AutoUpdate:  graphql.Boolean(plan.AutoUpdate.ValueBool()),
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create policy collection, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.AddPolicyCollection.Collection.ID)
	plan.UUID = types.StringValue(mutation.AddPolicyCollection.Collection.UUID)
	plan.Name = types.StringValue(mutation.AddPolicyCollection.Collection.Name)
	plan.Description = types.StringValue(mutation.AddPolicyCollection.Collection.Description)
	plan.CloudProvider = types.StringValue(mutation.AddPolicyCollection.Collection.Provider)
	plan.AutoUpdate = types.BoolValue(mutation.AddPolicyCollection.Collection.AutoUpdate)
	plan.System = types.BoolValue(mutation.AddPolicyCollection.Collection.System)
	plan.Repository = types.StringValue(mutation.AddPolicyCollection.Collection.Repository)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyCollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		PolicyCollection struct {
			ID          string
			UUID        string
			Name        string
			Description string
			Provider    string
			AutoUpdate  bool
			System      bool
			Repository  string
		} `graphql:"policyCollection(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(state.UUID.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy collection, got error: %s", err))
		return
	}

	if query.PolicyCollection.UUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(query.PolicyCollection.ID)
	state.UUID = types.StringValue(query.PolicyCollection.UUID)
	state.Name = types.StringValue(query.PolicyCollection.Name)
	state.Description = types.StringValue(query.PolicyCollection.Description)
	state.CloudProvider = types.StringValue(query.PolicyCollection.Provider)
	state.AutoUpdate = types.BoolValue(query.PolicyCollection.AutoUpdate)
	state.System = types.BoolValue(query.PolicyCollection.System)
	state.Repository = types.StringValue(query.PolicyCollection.Repository)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyCollectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		UpdatePolicyCollection struct {
			Collection struct {
				ID          string
				UUID        string
				Name        string
				Description string
				Provider    string
				AutoUpdate  bool
				System      bool
				Repository  string
			}
		} `graphql:"updatePolicyCollection(input: $input)"`
	}

	variables := map[string]any{
		"input": UpdatePolicyCollectionInput{
			UUID:        plan.UUID.ValueString(),
			Name:        graphql.String(plan.Name.ValueString()),
			Provider:    graphql.String(plan.CloudProvider.ValueString()),
			Description: graphql.String(plan.Description.ValueString()),
			AutoUpdate:  graphql.Boolean(plan.AutoUpdate.ValueBool()),
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update policy collection, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.UpdatePolicyCollection.Collection.ID)
	plan.UUID = types.StringValue(mutation.UpdatePolicyCollection.Collection.UUID)
	plan.Name = types.StringValue(mutation.UpdatePolicyCollection.Collection.Name)
	plan.Description = types.StringValue(mutation.UpdatePolicyCollection.Collection.Description)
	plan.CloudProvider = types.StringValue(mutation.UpdatePolicyCollection.Collection.Provider)
	plan.AutoUpdate = types.BoolValue(mutation.UpdatePolicyCollection.Collection.AutoUpdate)
	plan.System = types.BoolValue(mutation.UpdatePolicyCollection.Collection.System)
	plan.Repository = types.StringValue(mutation.UpdatePolicyCollection.Collection.Repository)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyCollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		RemovePolicyCollection struct {
			Collection struct {
				UUID string
			}
		} `graphql:"removePolicyCollection(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(state.UUID.ValueString()),
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete policy collection, got error: %s", err))
		return
	}
}

func (r *policyCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// GraphQL query to get policy collection by UUID
	var query struct {
		PolicyCollection struct {
			ID          string
			UUID        string
			Name        string
			Description string
			Provider    string
			AutoUpdate  bool
			System      bool
			Repository  string
		} `graphql:"policyCollection(uuid: $uuid)"`
	}

	variables := map[string]any{
		"uuid": graphql.String(req.ID),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Policy Collection",
			fmt.Sprintf("Could not read policy collection with UUID %s: %s", req.ID, err),
		)
		return
	}

	if query.PolicyCollection.UUID == "" {
		resp.Diagnostics.AddError(
			"Policy Collection Not Found",
			fmt.Sprintf("No policy collection found with UUID %s", req.ID),
		)
		return
	}

	// Set the imported attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), query.PolicyCollection.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), query.PolicyCollection.UUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), query.PolicyCollection.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), query.PolicyCollection.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), query.PolicyCollection.Provider)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("auto_update"), query.PolicyCollection.AutoUpdate)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("system"), query.PolicyCollection.System)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repository"), query.PolicyCollection.Repository)...)
}

type AddPolicyCollectionInput struct {
	Name        string          `json:"name"`
	Provider    string          `json:"provider"`
	Description graphql.String  `json:"description,omitempty"`
	AutoUpdate  graphql.Boolean `json:"autoUpdate,omitempty"`
}

type UpdatePolicyCollectionInput struct {
	UUID        string          `json:"uuid"`
	Name        graphql.String  `json:"name,omitempty"`
	Provider    graphql.String  `json:"provider,omitempty"`
	Description graphql.String  `json:"description,omitempty"`
	AutoUpdate  graphql.Boolean `json:"autoUpdate,omitempty"`
}
