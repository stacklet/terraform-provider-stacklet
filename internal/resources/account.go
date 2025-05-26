package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
	"github.com/stacklet/terraform-provider-stacklet/schemavalidate"
)

var (
	_ resource.Resource                = &accountResource{}
	_ resource.ResourceWithConfigure   = &accountResource{}
	_ resource.ResourceWithImportState = &accountResource{}
)

func NewAccountResource() resource.Resource {
	return &accountResource{}
}

type accountResource struct {
	api *api.API
}

func (r *accountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (r *accountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a Stacklet account with a specific cloud provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The cloud specific identifier for the account (e.g., AWS account ID, GCP project ID, Azure subscription UUID).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name for the account.",
				Required:    true,
			},
			"short_name": schema.StringAttribute{
				Description: "The short name for the account.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "More detailed information about the account.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
				Validators: []validator.String{
					schemavalidate.OneOfCloudProviders(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				Description: "The path used to group accounts in a hierarchy.",
				Optional:    true,
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email contact address for the account.",
				Optional:    true,
			},
			"security_context_wo": schema.StringAttribute{
				Description: "The input value for the security context for the account.",
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"security_context_wo_version": schema.StringAttribute{
				Description: "The version for the security context. Must be changed to update security_context_wo.",
				Optional:    true,
			},
			"security_context": schema.StringAttribute{
				Description: "The security context for the account, as returned by the API.",
				Computed:    true,
			},

			"variables": schema.StringAttribute{
				Description: "JSON encoded dict of values used for policy templating.",
				Optional:    true,
			},
		},
	}
}

func (r *accountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.AccountResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.AccountCreateInput{
		Name:            plan.Name.ValueString(),
		Key:             plan.Key.ValueString(),
		Provider:        api.CloudProvider(plan.CloudProvider.ValueString()),
		ShortName:       plan.ShortName.ValueStringPointer(),
		Description:     plan.Description.ValueStringPointer(),
		Email:           plan.Email.ValueStringPointer(),
		SecurityContext: config.SecurityContextWO.ValueStringPointer(),
		Variables:       plan.Variables.ValueStringPointer(),
	}
	account, err := r.api.Account.Create(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(updateAccountModel(&plan, account))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccountResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := r.api.Account.Read(ctx, state.CloudProvider.ValueString(), state.Key.ValueString())
	if err != nil {
		helpers.HandleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(updateAccountModel(&state, account))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.AccountResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var securityContext *string
	if state.SecurityContextWOVersion != plan.SecurityContextWOVersion {
		securityContext = config.SecurityContextWO.ValueStringPointer()
	}

	input := api.AccountUpdateInput{
		Key:             plan.Key.ValueString(),
		Provider:        api.CloudProvider(plan.CloudProvider.ValueString()),
		Name:            plan.Name.ValueStringPointer(),
		ShortName:       plan.ShortName.ValueStringPointer(),
		Description:     plan.Description.ValueStringPointer(),
		Email:           plan.Email.ValueStringPointer(),
		SecurityContext: securityContext,
		Variables:       plan.Variables.ValueStringPointer(),
	}

	account, err := r.api.Account.Update(ctx, input)
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(updateAccountModel(&plan, account))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AccountResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.Account.Delete(ctx, state.CloudProvider.ValueString(), state.Key.ValueString()); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := helpers.SplitImportID(req.ID, []string{"provider", "key"})
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), parts[1])...)
}

func updateAccountModel(m *models.AccountResource, account api.Account) diag.Diagnostic {
	m.ID = types.StringValue(account.ID)
	m.Key = types.StringValue(account.Key)
	m.Name = types.StringValue(account.Name)
	m.ShortName = tftypes.NullableString(account.ShortName)
	m.CloudProvider = types.StringValue(string(account.Provider))
	m.Description = tftypes.NullableString(account.Description)
	m.Path = tftypes.NullableString(account.Path)
	m.Email = tftypes.NullableString(account.Email)
	m.SecurityContext = tftypes.NullableString(account.SecurityContext)
	variablesString, err := tftypes.JSONString(account.Variables)
	if err != nil {
		return diag.NewErrorDiagnostic("Invalid content for variables", err.Error())
	}
	m.Variables = variablesString
	return nil
}
