package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ datasource.DataSource = &accountDataSource{}
)

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSource struct {
	client *graphql.Client
}

type accountDataSourceModel struct {
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
	Variables       types.Map    `tfsdk:"variables"`
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a cloud account in Stacklet across different cloud providers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "The cloud specific identifier for the account (e.g., AWS account ID, GCP project ID, Azure subscription UUID).",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The human readable identifier for the account.",
				Computed:    true,
			},
			"short_name": schema.StringAttribute{
				Description: "The short name for the account.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "More detailed information about the account.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account (aws, azure, gcp, kubernetes, or tencentcloud).",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "The path used to group accounts in a hierarchy.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email contact address for the account.",
				Computed:    true,
			},
			"security_context": schema.StringAttribute{
				Description: "The security context for the account.",
				Computed:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether the account is active or has been deactivated.",
				Computed:    true,
			},
			"variables": schema.MapAttribute{
				ElementType: types.StringType,
				Description: "Values used for policy templating.",
				Computed:    true,
			},
		},
	}
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
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

	provider, err := NewCloudProvider(data.CloudProvider.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Provider", err.Error())
		return
	}

	variables := map[string]any{
		"provider": provider,
		"key":      graphql.String(data.Key.ValueString()),
	}

	err = d.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account, got error: %s", err))
		return
	}

	if query.Account.Key == "" {
		resp.Diagnostics.AddError("Not Found", "No account found with the specified provider and key")
		return
	}

	data.ID = types.StringValue(query.Account.ID)
	data.Key = types.StringValue(query.Account.Key)
	data.Name = types.StringValue(query.Account.Name)
	data.ShortName = nullableString(query.Account.ShortName)
	data.Description = nullableString(query.Account.Description)
	data.CloudProvider = types.StringValue(string(query.Account.Provider))
	data.Path = nullableString(query.Account.Path)
	data.Email = nullableString(query.Account.Email)
	data.SecurityContext = nullableString(query.Account.SecurityContext)
	data.Active = types.BoolValue(query.Account.Active)
	if vars, diags := jsonMap(ctx, query.Account.Variables); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	} else {
		data.Variables = vars
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
