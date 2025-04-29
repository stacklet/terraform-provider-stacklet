package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				Provider        CloudProvider
				Path            string
				Email           string
				SecurityContext string
				Active          bool
				Variables       string
			}
		} `graphql:"addAccount(input: $input)"`
	}

	// Convert provider string to CloudProvider type
	provider := CloudProvider(strings.ToUpper(plan.CloudProvider.ValueString()))
	if err := provider.Validate(); err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	input := map[string]any{
		"input": AccountInput{
			Name:     plan.Name.ValueString(),
			Key:      plan.Key.ValueString(),
			Provider: provider,
			ShortName: func() *string {
				// Always send empty string if not explicitly set
				if !plan.ShortName.IsNull() {
					s := plan.ShortName.ValueString()
					return &s
				}
				s := ""
				return &s
			}(),
			Description: func() *string {
				if !plan.Description.IsNull() {
					s := plan.Description.ValueString()
					return &s
				}
				return nil
			}(),
			Path: func() *string {
				if !plan.Path.IsNull() {
					s := plan.Path.ValueString()
					return &s
				}
				return nil
			}(),
			Email: func() *string {
				if !plan.Email.IsNull() {
					s := plan.Email.ValueString()
					return &s
				}
				return nil
			}(),
			SecurityContext: func() *string {
				if !plan.SecurityContext.IsNull() {
					s := plan.SecurityContext.ValueString()
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
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create account, got error: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(mutation.AddAccount.Account.ID)
	plan.Key = types.StringValue(mutation.AddAccount.Account.Key)
	plan.Name = types.StringValue(mutation.AddAccount.Account.Name)
	plan.CloudProvider = types.StringValue(string(mutation.AddAccount.Account.Provider))
	plan.ShortName = types.StringValue(mutation.AddAccount.Account.ShortName)

	if mutation.AddAccount.Account.Description != "" {
		plan.Description = types.StringValue(mutation.AddAccount.Account.Description)
	} else if !plan.Description.IsNull() {
		// Keep the planned value
	} else {
		plan.Description = types.StringNull()
	}

	if mutation.AddAccount.Account.Path != "" {
		plan.Path = types.StringValue(mutation.AddAccount.Account.Path)
	} else if !plan.Path.IsNull() {
		// Keep the planned value
	} else {
		plan.Path = types.StringNull()
	}

	if mutation.AddAccount.Account.Email != "" {
		plan.Email = types.StringValue(mutation.AddAccount.Account.Email)
	} else if !plan.Email.IsNull() {
		// Keep the planned value
	} else {
		plan.Email = types.StringNull()
	}

	if mutation.AddAccount.Account.SecurityContext != "" {
		plan.SecurityContext = types.StringValue(mutation.AddAccount.Account.SecurityContext)
	} else if !plan.SecurityContext.IsNull() {
		// Keep the planned value
	} else {
		plan.SecurityContext = types.StringNull()
	}

	if !plan.Active.IsNull() {
		// Keep the planned value
	} else {
		plan.Active = types.BoolValue(mutation.AddAccount.Account.Active)
	}

	if mutation.AddAccount.Account.Variables != "" {
		plan.Variables = types.StringValue(mutation.AddAccount.Account.Variables)
	} else if !plan.Variables.IsNull() {
		// Keep the planned value
	} else {
		plan.Variables = types.StringNull()
	}

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
			Provider        CloudProvider
			Path            string
			Email           string
			SecurityContext string
			Active          bool
			Variables       string
		} `graphql:"account(provider: $provider, key: $key)"`
	}

	// Convert provider string to CloudProvider type
	provider := CloudProvider(strings.ToUpper(state.CloudProvider.ValueString()))
	if err := provider.Validate(); err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	variables := map[string]any{
		"provider": provider,
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
	state.CloudProvider = types.StringValue(string(query.Account.Provider))

	if query.Account.Description == "" {
		state.Description = types.StringNull()
	} else {
		state.Description = types.StringValue(query.Account.Description)
	}

	if query.Account.Path == "" {
		state.Path = types.StringNull()
	} else {
		state.Path = types.StringValue(query.Account.Path)
	}

	if query.Account.Email == "" {
		state.Email = types.StringNull()
	} else {
		state.Email = types.StringValue(query.Account.Email)
	}

	if query.Account.SecurityContext == "" {
		state.SecurityContext = types.StringNull()
	} else {
		state.SecurityContext = types.StringValue(query.Account.SecurityContext)
	}

	state.Active = types.BoolValue(query.Account.Active)

	if query.Account.Variables == "" {
		state.Variables = types.StringNull()
	} else {
		state.Variables = types.StringValue(query.Account.Variables)
	}

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
				Provider        CloudProvider
				Path            string
				Email           string
				SecurityContext string
				Active          bool
				Variables       string
			}
		} `graphql:"updateAccount(input: $input)"`
	}

	// Convert provider string to CloudProvider type
	provider := CloudProvider(strings.ToUpper(plan.CloudProvider.ValueString()))
	if err := provider.Validate(); err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	input := map[string]any{
		"input": UpdateAccountInput{
			Key:      plan.Key.ValueString(),
			Provider: provider,
			Name:     plan.Name.ValueString(),
			ShortName: func() *string {
				// Always send empty string if not explicitly set
				if !plan.ShortName.IsNull() {
					s := plan.ShortName.ValueString()
					return &s
				}
				s := ""
				return &s
			}(),
			Description: func() *string {
				if !plan.Description.IsNull() {
					s := plan.Description.ValueString()
					return &s
				}
				return nil
			}(),
			Email: func() *string {
				if !plan.Email.IsNull() {
					s := plan.Email.ValueString()
					return &s
				}
				return nil
			}(),
			SecurityContext: func() *string {
				if !plan.SecurityContext.IsNull() {
					s := plan.SecurityContext.ValueString()
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update account, got error: %s", err))
		return
	}

	plan.ID = types.StringValue(mutation.UpdateAccount.Account.ID)
	plan.Key = types.StringValue(mutation.UpdateAccount.Account.Key)
	plan.Name = types.StringValue(mutation.UpdateAccount.Account.Name)
	plan.CloudProvider = types.StringValue(string(mutation.UpdateAccount.Account.Provider))
	plan.ShortName = types.StringValue(mutation.UpdateAccount.Account.ShortName)

	if mutation.UpdateAccount.Account.Description != "" {
		plan.Description = types.StringValue(mutation.UpdateAccount.Account.Description)
	} else if !plan.Description.IsNull() {
		// Keep the planned value
	} else {
		plan.Description = types.StringNull()
	}

	if mutation.UpdateAccount.Account.Path != "" {
		plan.Path = types.StringValue(mutation.UpdateAccount.Account.Path)
	} else if !plan.Path.IsNull() {
		// Keep the planned value
	} else {
		plan.Path = types.StringNull()
	}

	if mutation.UpdateAccount.Account.Email != "" {
		plan.Email = types.StringValue(mutation.UpdateAccount.Account.Email)
	} else if !plan.Email.IsNull() {
		// Keep the planned value
	} else {
		plan.Email = types.StringNull()
	}

	if mutation.UpdateAccount.Account.SecurityContext != "" {
		plan.SecurityContext = types.StringValue(mutation.UpdateAccount.Account.SecurityContext)
	} else if !plan.SecurityContext.IsNull() {
		// Keep the planned value
	} else {
		plan.SecurityContext = types.StringNull()
	}

	if !plan.Active.IsNull() {
		// Keep the planned value
	} else {
		plan.Active = types.BoolValue(mutation.UpdateAccount.Account.Active)
	}

	if mutation.UpdateAccount.Account.Variables != "" {
		plan.Variables = types.StringValue(mutation.UpdateAccount.Account.Variables)
	} else if !plan.Variables.IsNull() {
		// Keep the planned value
	} else {
		plan.Variables = types.StringNull()
	}

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
			}
		} `graphql:"removeAccount(provider: $provider, key: $key)"`
	}

	// Convert provider string to CloudProvider type
	provider := CloudProvider(strings.ToUpper(state.CloudProvider.ValueString()))
	if err := provider.Validate(); err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	variables := map[string]any{
		"provider": provider,
		"key":      graphql.String(state.Key.ValueString()),
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

type AccountInput struct {
	Name            string        `json:"name"`
	Key             string        `json:"key"`
	Provider        CloudProvider `json:"provider"`
	ShortName       *string       `json:"shortName,omitempty"`
	Description     *string       `json:"description,omitempty"`
	Path            *string       `json:"path,omitempty"`
	Email           *string       `json:"email,omitempty"`
	SecurityContext *string       `json:"securityContext,omitempty"`
	Variables       *string       `json:"variables,omitempty"`
}

type UpdateAccountInput struct {
	Key             string        `json:"key"`
	Provider        CloudProvider `json:"provider"`
	Name            string        `json:"name"`
	ShortName       *string       `json:"shortName,omitempty"`
	Description     *string       `json:"description,omitempty"`
	Email           *string       `json:"email,omitempty"`
	SecurityContext *string       `json:"securityContext,omitempty"`
	Variables       *string       `json:"variables,omitempty"`
}
