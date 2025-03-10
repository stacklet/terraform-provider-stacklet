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

var (
	_ resource.Resource                = &bindingResource{}
	_ resource.ResourceWithImportState = &bindingResource{}
)

func NewBindingResource() resource.Resource {
	return &bindingResource{}
}

type bindingResource struct {
	client *graphql.Client
}

type bindingResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	Schedule             types.String `tfsdk:"schedule"`
	Variables            types.String `tfsdk:"variables"`
	AccountGroupUUID     types.String `tfsdk:"account_group_uuid"`
	PolicyCollectionUUID types.String `tfsdk:"policy_collection_uuid"`
	Deploy               types.Bool   `tfsdk:"deploy"`
}

func (r *bindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding"
}

func (r *bindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a binding between an account group and a policy collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the binding.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the binding.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the binding.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the binding.",
				Optional:    true,
			},
			"auto_deploy": schema.BoolAttribute{
				Description: "Whether the binding should automatically deploy when the policy collection changes.",
				Optional:    true,
			},
			"schedule": schema.StringAttribute{
				Description: "The schedule for the binding (e.g., 'rate(1 hour)', 'rate(2 hours)', or cron expression).",
				Optional:    true,
			},
			"variables": schema.StringAttribute{
				Description: "JSON-encoded dictionary of values used for policy templating.",
				Optional:    true,
			},
			"account_group_uuid": schema.StringAttribute{
				Description: "The UUID of the account group this binding applies to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_collection_uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection this binding applies.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"deploy": schema.BoolAttribute{
				Description: "Whether to deploy the binding immediately after creation.",
				Optional:    true,
			},
		},
	}
}

func (r *bindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		AddBinding struct {
			Binding struct {
				ID   string
				UUID string
			}
		} `graphql:"addBinding(input: $input)"`
	}

	input := map[string]interface{}{
		"input": AddBindingInput{
			Name: plan.Name.ValueString(),
			Description: func() *string {
				if !plan.Description.IsNull() {
					s := plan.Description.ValueString()
					return &s
				}
				return nil
			}(),
			AutoDeploy: func() *bool {
				if !plan.AutoDeploy.IsNull() {
					b := plan.AutoDeploy.ValueBool()
					return &b
				}
				return nil
			}(),
			Schedule: func() *string {
				if !plan.Schedule.IsNull() {
					s := plan.Schedule.ValueString()
					return &s
				}
				return nil
			}(),
			Variables: func() *string {
				if !plan.Variables.IsNull() {
					s := plan.Variables.ValueString()
					return &s
				}
				return nil
			}(),
			AccountGroupUUID:     plan.AccountGroupUUID.ValueString(),
			PolicyCollectionUUID: plan.PolicyCollectionUUID.ValueString(),
			Deploy: func() *bool {
				if !plan.Deploy.IsNull() {
					b := plan.Deploy.ValueBool()
					return &b
				}
				return nil
			}(),
		},
	}

	err := r.client.Mutate(ctx, &mutation, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create binding, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.AddBinding.Binding.ID)
	plan.UUID = types.StringValue(mutation.AddBinding.Binding.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query struct {
		Binding struct {
			ID           string
			UUID         string
			Name         string
			Description  string
			AutoDeploy   bool
			Schedule     string
			Variables    string
			AccountGroup struct {
				UUID string
			}
			PolicyCollection struct {
				UUID string
			}
		} `graphql:"binding(uuid: $uuid)"`
	}

	variables := map[string]interface{}{
		"uuid": graphql.String(state.UUID.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read binding, got error: %s", err))
		return
	}

	if query.Binding.UUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(query.Binding.ID)
	state.UUID = types.StringValue(query.Binding.UUID)
	state.Name = types.StringValue(query.Binding.Name)
	state.Description = types.StringValue(query.Binding.Description)
	state.AutoDeploy = types.BoolValue(query.Binding.AutoDeploy)
	state.Schedule = types.StringValue(query.Binding.Schedule)
	state.Variables = types.StringValue(query.Binding.Variables)
	state.AccountGroupUUID = types.StringValue(query.Binding.AccountGroup.UUID)
	state.PolicyCollectionUUID = types.StringValue(query.Binding.PolicyCollection.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan bindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		UpdateBinding struct {
			Binding struct {
				ID   string
				UUID string
			}
		} `graphql:"updateBinding(input: $input)"`
	}

	input := map[string]interface{}{
		"input": UpdateBindingInput{
			UUID: plan.UUID.ValueString(),
			Name: plan.Name.ValueString(),
			Description: func() *string {
				if !plan.Description.IsNull() {
					s := plan.Description.ValueString()
					return &s
				}
				return nil
			}(),
			AutoDeploy: func() *bool {
				if !plan.AutoDeploy.IsNull() {
					b := plan.AutoDeploy.ValueBool()
					return &b
				}
				return nil
			}(),
			Schedule: func() *string {
				if !plan.Schedule.IsNull() {
					s := plan.Schedule.ValueString()
					return &s
				}
				return nil
			}(),
			Variables: func() *string {
				if !plan.Variables.IsNull() {
					s := plan.Variables.ValueString()
					return &s
				}
				return nil
			}(),
		},
	}

	err := r.client.Mutate(ctx, &mutation, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update binding, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.UpdateBinding.Binding.ID)
	plan.UUID = types.StringValue(mutation.UpdateBinding.Binding.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		RemoveBinding struct {
			Binding struct {
				UUID string
			}
		} `graphql:"removeBinding(uuid: $uuid)"`
	}

	variables := map[string]interface{}{
		"uuid": graphql.String(state.UUID.ValueString()),
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete binding, got error: %s", err))
		return
	}
}

func (r *bindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), req.ID)...)
}

type AddBindingInput struct {
	Name                 string  `json:"name"`
	Description          *string `json:"description,omitempty"`
	AutoDeploy           *bool   `json:"autoDeploy,omitempty"`
	Schedule             *string `json:"schedule,omitempty"`
	Variables            *string `json:"variables,omitempty"`
	AccountGroupUUID     string  `json:"accountGroupUUID"`
	PolicyCollectionUUID string  `json:"policyCollectionUUID"`
	Deploy               *bool   `json:"deploy,omitempty"`
}

type UpdateBindingInput struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	AutoDeploy  *bool   `json:"autoDeploy,omitempty"`
	Schedule    *string `json:"schedule,omitempty"`
	Variables   *string `json:"variables,omitempty"`
}
