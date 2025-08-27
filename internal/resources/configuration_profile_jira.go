// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ resource.Resource                = &configurationProfileJiraResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileJiraResource{}
	_ resource.ResourceWithImportState = &configurationProfileJiraResource{}
)

func NewConfigurationProfileJiraResource() resource.Resource {
	return &configurationProfileJiraResource{}
}

type configurationProfileJiraResource struct {
	api *api.API
}

func (r *configurationProfileJiraResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_jira"
}

func (r *configurationProfileJiraResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the Jira configuration profile.

The profile is global, adding multiple resources of this kind will cause them to override each other.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The Jira instance URL.",
				Optional:    true,
			},
			"user": schema.StringAttribute{
				Description: "The Jira instance authentication username.",
				Required:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The encrypted value for the API key, returned by the API.",
				Computed:    true,
			},
			"api_key_wo": schema.StringAttribute{
				Description: "The input value for the API key.",
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"api_key_wo_version": schema.StringAttribute{
				Description: "The version for the API key. Must be changed to update api_key_wo.",
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"project": schema.ListNestedBlock{
				Description: "Jira project configuration.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"closed_status": schema.StringAttribute{
							Description: "The state for closed tickets.",
							Required:    true,
						},
						"issue_type": schema.StringAttribute{
							Description: "The type of issue to use for tickets.",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the project.",
							Required:    true,
						},
						"project": schema.StringAttribute{
							Description: "The ID of the project.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *configurationProfileJiraResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileJiraResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ConfigurationProfileJiraResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.api.ConfigurationProfile.ReadJira(ctx)
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateJiraModel(&state, config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configurationProfileJiraResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.ConfigurationProfileJiraResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects, diags := r.getProjects(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.JiraConfiguration{
		URL:      plan.URL.ValueStringPointer(),
		Projects: projects,
		User:     plan.User.ValueString(),
		APIKey:   config.APIKeyWO.ValueString(),
	}
	jiraConfig, err := r.api.ConfigurationProfile.UpsertJira(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateJiraModel(&plan, jiraConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileJiraResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.ConfigurationProfileJiraResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiKey string
	if state.APIKeyWOVersion == plan.APIKeyWOVersion {
		apiKey = state.APIKey.ValueString() // send back the previous encrypted value
	} else {
		apiKey = config.APIKeyWO.ValueString() // send the new value from the config
	}

	projects, diags := r.getProjects(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.JiraConfiguration{
		URL:      plan.URL.ValueStringPointer(),
		Projects: projects,
		User:     plan.User.ValueString(),
		APIKey:   apiKey,
	}
	profileConfig, err := r.api.ConfigurationProfile.UpsertJira(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateJiraModel(&plan, profileConfig)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configurationProfileJiraResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileJiraResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteJira(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileJiraResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileJira))...)
}

func (r configurationProfileJiraResource) updateJiraModel(m *models.ConfigurationProfileJiraResource, config *api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	jiraConfig := config.Record.JiraConfiguration

	m.ID = types.StringValue(config.ID)
	m.Profile = types.StringValue(config.Profile)
	m.URL = types.StringPointerValue(jiraConfig.URL)
	m.User = types.StringValue(jiraConfig.User)
	m.APIKey = types.StringValue(jiraConfig.APIKey)

	// get the current names ordering to preserve it in the updated projects block
	names := models.ListItemsIdentifiers(m.Projects, "name")

	updater := modelupdate.NewConfigurationProfileUpdater(*config)
	projects, diags := updater.JiraProjects(names)
	if diags.HasError() {
		return diags
	}
	m.Projects = projects
	return diags
}

func (r configurationProfileJiraResource) getProjects(ctx context.Context, m models.ConfigurationProfileJiraResource) ([]api.JiraProject, diag.Diagnostics) {
	if m.Projects.IsNull() {
		return nil, nil
	}

	var diags diag.Diagnostics

	projects := []api.JiraProject{}
	for i, elem := range m.Projects.Elements() {
		block, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddAttributeError(
				path.Root(fmt.Sprintf("project.%d", i)),
				"Invalid project configuration",
				"Project configuration block is invalid.",
			)
			return nil, diags
		}
		var p models.JiraProject
		if diags := block.As(ctx, &p, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}

		projects = append(
			projects,
			api.JiraProject{
				ClosedStatus: p.ClosedStatus.ValueString(),
				IssueType:    p.IssueType.ValueString(),
				Name:         p.Name.ValueString(),
				Project:      p.Project.ValueString(),
			},
		)
	}
	return projects, diags
}
