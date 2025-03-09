package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ resource.Resource = &accountGroupResource{}
)

func NewAccountGroupResource() resource.Resource {
	return &accountGroupResource{}
}

type accountGroupResource struct {
	client *graphql.Client
}

type accountGroupResourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
}

func (r *accountGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group"
}

func (r *accountGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an account group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the account group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the account group.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the account group.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account group (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
			},
		},
	}
}

func (r *accountGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *accountGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accountGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		AddAccountGroup struct {
			Group struct {
				ID          string
				UUID        string
				Name        string
				Description string
				Provider    string
			}
		} `graphql:"addAccountGroup(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": AddAccountGroupInput{
			Name:        plan.Name.ValueString(),
			Provider:    plan.CloudProvider.ValueString(),
			Description: graphql.String(plan.Description.ValueString()),
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account group, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.AddAccountGroup.Group.ID)
	plan.UUID = types.StringValue(mutation.AddAccountGroup.Group.UUID)
	plan.Name = types.StringValue(mutation.AddAccountGroup.Group.Name)
	plan.Description = types.StringValue(mutation.AddAccountGroup.Group.Description)
	plan.CloudProvider = types.StringValue(mutation.AddAccountGroup.Group.Provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		AccountGroup struct {
			ID          string
			UUID        string
			Name        string
			Description string
			Provider    string
		} `graphql:"accountGroup(uuid: $uuid)"`
	}

	variables := map[string]interface{}{
		"uuid": graphql.String(state.UUID.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account group, got error: %s", err))
		return
	}

	if query.AccountGroup.UUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(query.AccountGroup.ID)
	state.UUID = types.StringValue(query.AccountGroup.UUID)
	state.Name = types.StringValue(query.AccountGroup.Name)
	state.Description = types.StringValue(query.AccountGroup.Description)
	state.CloudProvider = types.StringValue(query.AccountGroup.Provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accountGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		UpdateAccountGroup struct {
			Group struct {
				ID          string
				UUID        string
				Name        string
				Description string
				Provider    string
			}
		} `graphql:"updateAccountGroup(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": UpdateAccountGroupInput{
			UUID:        plan.UUID.ValueString(),
			Name:        graphql.String(plan.Name.ValueString()),
			Provider:    graphql.String(plan.CloudProvider.ValueString()),
			Description: graphql.String(plan.Description.ValueString()),
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update account group, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.UpdateAccountGroup.Group.ID)
	plan.UUID = types.StringValue(mutation.UpdateAccountGroup.Group.UUID)
	plan.Name = types.StringValue(mutation.UpdateAccountGroup.Group.Name)
	plan.Description = types.StringValue(mutation.UpdateAccountGroup.Group.Description)
	plan.CloudProvider = types.StringValue(mutation.UpdateAccountGroup.Group.Provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accountGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL mutation
	var mutation struct {
		RemoveAccountGroup struct {
			Group struct {
				UUID string
			}
		} `graphql:"removeAccountGroup(uuid: $uuid)"`
	}

	variables := map[string]interface{}{
		"uuid": graphql.String(state.UUID.ValueString()),
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete account group, got error: %s", err))
		return
	}
}

type AddAccountGroupInput struct {
	Name        string         `json:"name"`
	Provider    string         `json:"provider"`
	Description graphql.String `json:"description,omitempty"`
}

type UpdateAccountGroupInput struct {
	UUID        string         `json:"uuid"`
	Name        graphql.String `json:"name,omitempty"`
	Provider    graphql.String `json:"provider,omitempty"`
	Description graphql.String `json:"description,omitempty"`
}
