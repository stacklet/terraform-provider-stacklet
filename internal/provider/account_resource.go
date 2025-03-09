package provider

import (
	"context"
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
	_ resource.Resource                = &accountResource{}
	_ resource.ResourceWithImportState = &accountResource{}
)

func NewAccountResource() resource.Resource {
	return &accountResource{}
}

type accountResource struct {
	client *graphql.Client
}

type accountResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Key             types.String `tfsdk:"key"`
	Name            types.String `tfsdk:"name"`
	ShortName       types.String `tfsdk:"short_name"`
	Description     types.String `tfsdk:"description"`
	CloudProvider   types.String `tfsdk:"cloud_provider"`
	Path            types.String `tfsdk:"path"`
	Email           types.String `tfsdk:"email"`
	SecurityContext types.String `tfsdk:"security_context"`
	Active          types.Bool   `tfsdk:"active"`
	Variables       types.String `tfsdk:"variables"`
}

func (r *accountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (r *accountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Stacklet account.",
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
				Description: "The human readable identifier for the account.",
				Required:    true,
			},
			"short_name": schema.StringAttribute{
				Description: "The short name used as a column header if set.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "More detailed information about the account.",
				Optional:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				Description: "The path used to group accounts in a hierarchy.",
				Optional:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email contact address for the account.",
				Optional:    true,
			},
			"security_context": schema.StringAttribute{
				Description: "The security context for the account.",
				Optional:    true,
				Sensitive:   true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether the account is active or has been deactivated.",
				Optional:    true,
			},
			"variables": schema.StringAttribute{
				Description: "JSON encoded dict of values used for policy templating.",
				Optional:    true,
			},
		},
	}
}

func (r *accountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		AddAccount struct {
			Account struct {
				ID              string
				Key             string
				Name            string
				ShortName       string
				Description     string
				Provider        string
				Path            string
				Email           string
				SecurityContext string
				Active          bool
				Variables       string
			} `graphql:"account"`
		} `graphql:"addAccount(input: $input)"`
	}

	input := map[string]interface{}{
		"key":      graphql.String(plan.Key.ValueString()),
		"name":     graphql.String(plan.Name.ValueString()),
		"provider": graphql.String(plan.CloudProvider.ValueString()),
	}

	if !plan.ShortName.IsNull() {
		input["shortName"] = graphql.String(plan.ShortName.ValueString())
	}
	if !plan.Description.IsNull() {
		input["description"] = graphql.String(plan.Description.ValueString())
	}
	if !plan.Path.IsNull() {
		input["path"] = graphql.String(plan.Path.ValueString())
	}
	if !plan.Email.IsNull() {
		input["email"] = graphql.String(plan.Email.ValueString())
	}
	if !plan.SecurityContext.IsNull() {
		input["securityContext"] = graphql.String(plan.SecurityContext.ValueString())
	}
	if !plan.Active.IsNull() {
		input["active"] = graphql.Boolean(plan.Active.ValueBool())
	}
	if !plan.Variables.IsNull() {
		input["variables"] = graphql.String(plan.Variables.ValueString())
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.AddAccount.Account.ID)
	plan.Key = types.StringValue(mutation.AddAccount.Account.Key)
	plan.Name = types.StringValue(mutation.AddAccount.Account.Name)
	plan.ShortName = types.StringValue(mutation.AddAccount.Account.ShortName)
	plan.Description = types.StringValue(mutation.AddAccount.Account.Description)
	plan.CloudProvider = types.StringValue(mutation.AddAccount.Account.Provider)
	plan.Path = types.StringValue(mutation.AddAccount.Account.Path)
	plan.Email = types.StringValue(mutation.AddAccount.Account.Email)
	plan.SecurityContext = types.StringValue(mutation.AddAccount.Account.SecurityContext)
	plan.Active = types.BoolValue(mutation.AddAccount.Account.Active)
	plan.Variables = types.StringValue(mutation.AddAccount.Account.Variables)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query struct {
		Account struct {
			ID              string
			Key             string
			Name            string
			ShortName       string
			Description     string
			Provider        string
			Path            string
			Email           string
			SecurityContext string
			Active          bool
			Variables       string
		} `graphql:"account(provider: $provider, key: $key)"`
	}

	variables := map[string]interface{}{
		"provider": graphql.String(state.CloudProvider.ValueString()),
		"key":      graphql.String(state.Key.ValueString()),
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account, got error: %s", err))
		return
	}

	if query.Account.Key == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(query.Account.ID)
	state.Key = types.StringValue(query.Account.Key)
	state.Name = types.StringValue(query.Account.Name)
	state.ShortName = types.StringValue(query.Account.ShortName)
	state.Description = types.StringValue(query.Account.Description)
	state.CloudProvider = types.StringValue(query.Account.Provider)
	state.Path = types.StringValue(query.Account.Path)
	state.Email = types.StringValue(query.Account.Email)
	state.SecurityContext = types.StringValue(query.Account.SecurityContext)
	state.Active = types.BoolValue(query.Account.Active)
	state.Variables = types.StringValue(query.Account.Variables)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		UpdateAccount struct {
			Account struct {
				ID              string
				Key             string
				Name            string
				ShortName       string
				Description     string
				Provider        string
				Path            string
				Email           string
				SecurityContext string
				Active          bool
				Variables       string
			} `graphql:"account"`
		} `graphql:"updateAccount(input: $input)"`
	}

	input := map[string]interface{}{
		"key":      graphql.String(plan.Key.ValueString()),
		"name":     graphql.String(plan.Name.ValueString()),
		"provider": graphql.String(plan.CloudProvider.ValueString()),
	}

	if !plan.ShortName.IsNull() {
		input["shortName"] = graphql.String(plan.ShortName.ValueString())
	}
	if !plan.Description.IsNull() {
		input["description"] = graphql.String(plan.Description.ValueString())
	}
	if !plan.Path.IsNull() {
		input["path"] = graphql.String(plan.Path.ValueString())
	}
	if !plan.Email.IsNull() {
		input["email"] = graphql.String(plan.Email.ValueString())
	}
	if !plan.SecurityContext.IsNull() {
		input["securityContext"] = graphql.String(plan.SecurityContext.ValueString())
	}
	if !plan.Active.IsNull() {
		input["active"] = graphql.Boolean(plan.Active.ValueBool())
	}
	if !plan.Variables.IsNull() {
		input["variables"] = graphql.String(plan.Variables.ValueString())
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update account, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.UpdateAccount.Account.ID)
	plan.Key = types.StringValue(mutation.UpdateAccount.Account.Key)
	plan.Name = types.StringValue(mutation.UpdateAccount.Account.Name)
	plan.ShortName = types.StringValue(mutation.UpdateAccount.Account.ShortName)
	plan.Description = types.StringValue(mutation.UpdateAccount.Account.Description)
	plan.CloudProvider = types.StringValue(mutation.UpdateAccount.Account.Provider)
	plan.Path = types.StringValue(mutation.UpdateAccount.Account.Path)
	plan.Email = types.StringValue(mutation.UpdateAccount.Account.Email)
	plan.SecurityContext = types.StringValue(mutation.UpdateAccount.Account.SecurityContext)
	plan.Active = types.BoolValue(mutation.UpdateAccount.Account.Active)
	plan.Variables = types.StringValue(mutation.UpdateAccount.Account.Variables)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mutation struct {
		RemoveAccount struct {
			Account struct {
				Key string
			} `graphql:"account"`
		} `graphql:"removeAccount(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"key":      graphql.String(state.Key.ValueString()),
			"provider": graphql.String(state.CloudProvider.ValueString()),
		},
	}

	err := r.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete account, got error: %s", err))
		return
	}
}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Extract the provider and key from the import ID (format: provider:key)
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: provider:key (e.g., aws:123456789012)",
		)
		return
	}

	provider := idParts[0]
	key := idParts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), provider)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
}
