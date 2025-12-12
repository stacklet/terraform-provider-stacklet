// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/planmodifiers"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemavalidate"
)

var (
	_ resource.Resource                = &policyCollectionResource{}
	_ resource.ResourceWithConfigure   = &policyCollectionResource{}
	_ resource.ResourceWithImportState = &policyCollectionResource{}
)

func newPolicyCollectionResource() resource.Resource {
	return &policyCollectionResource{}
}

type policyCollectionResource struct {
	apiResource
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
				Validators: []validator.String{
					schemavalidate.OneOfCloudProviders(),
				},
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
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dynamic": schema.BoolAttribute{
				Description: "Whether this is a dynamic policy collection.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dynamic_config": schema.SingleNestedAttribute{
				Description: "Configuration for dynamic behavior.",
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					planmodifiers.RequiresReplaceIfFieldsChanged("repository_uuid"),
				},
				Attributes: map[string]schema.Attribute{
					"repository_uuid": schema.StringAttribute{
						Description: "The UUID of the repository the collection is linked to.",
						Required:    true,
					},
					"namespace": schema.StringAttribute{
						Description: "The namespace for policies from the repository.",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"branch_name": schema.StringAttribute{
						Description: "The repository branch from which policies are imported.",
						Optional:    true,
						Computed:    true,
					},
					"policy_directories": schema.ListAttribute{
						Description: "Optional list of subdirectory to limit the scan to.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"policy_file_suffixes": schema.ListAttribute{
						Description: "Optional list of suffixes for policy files to limit the scan to.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"role_assignment_target": schema.StringAttribute{
				Description: "An opaque identifier for role assignments. Use this value when assigning roles to this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *policyCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PolicyCollectionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryUUID, repositoryView, diags := r.getDynamicDetails(ctx, plan.DynamicConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionCreateInput{
		Name:        plan.Name.ValueString(),
		Provider:    api.CloudProvider(plan.CloudProvider.ValueString()),
		Description: plan.Description.ValueStringPointer(),
		AutoUpdate:  plan.AutoUpdate.ValueBoolPointer(),
		// the following are null if the policy collection is not dynamic
		RepositoryUUID: repositoryUUID,
		RepositoryView: repositoryView,
	}
	policyCollection, err := r.api.PolicyCollection.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, policyCollection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PolicyCollectionResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyCollection, err := r.api.PolicyCollection.Read(ctx, state.UUID.ValueString(), "")
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(ctx, policyCollection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan models.PolicyCollectionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, repositoryView, diags := r.getDynamicDetails(ctx, plan.DynamicConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.PolicyCollectionUpdateInput{
		UUID:           plan.UUID.ValueString(),
		Name:           plan.Name.ValueString(),
		Provider:       api.CloudProvider(plan.CloudProvider.ValueString()),
		Description:    plan.Description.ValueStringPointer(),
		AutoUpdate:     plan.AutoUpdate.ValueBoolPointer(),
		RepositoryView: repositoryView, // null if the policy collection is not dynamic
	}

	policyCollection, err := r.api.PolicyCollection.Update(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, policyCollection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PolicyCollectionResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.PolicyCollection.Delete(ctx, state.UUID.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *policyCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"uuid"})
}

func (r policyCollectionResource) getDynamicDetails(ctx context.Context, planDynamicConfig types.Object) (*string, *api.RepositoryViewInput, diag.Diagnostics) {
	var uuid *string
	var view *api.RepositoryViewInput
	var diags diag.Diagnostics

	if !planDynamicConfig.IsNull() {
		var dynamicConfig models.PolicyCollectionDynamicConfig
		diags = planDynamicConfig.As(ctx, &dynamicConfig, objectAsOptions)
		uuid = dynamicConfig.RepositoryUUID.ValueStringPointer()
		view = &api.RepositoryViewInput{
			BranchName:        dynamicConfig.BranchName.ValueStringPointer(),
			PolicyDirectories: api.StringsList(dynamicConfig.PolicyDirectories),
			PolicyFileSuffix:  api.StringsList(dynamicConfig.PolicyFileSuffixes),
		}
	}

	return uuid, view, diags
}
