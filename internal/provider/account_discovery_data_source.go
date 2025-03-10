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
	_ datasource.DataSource = &accountDiscoveryDataSource{}
)

func NewAccountDiscoveryDataSource() datasource.DataSource {
	return &accountDiscoveryDataSource{}
}

type accountDiscoveryDataSource struct {
	client *graphql.Client
}

type accountDiscoveryDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Config        types.String `tfsdk:"config"`
	Schedule      types.String `tfsdk:"schedule"`
	Validity      types.String `tfsdk:"validity"`
}

func (d *accountDiscoveryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_discovery"
}

func (d *accountDiscoveryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch an account discovery configuration by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account discovery configuration.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the account discovery configuration.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the account discovery configuration.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Human-readable notes about the account discovery configuration.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Computed:    true,
				Description: "The cloud provider for the account discovery (aws, azure, gcp, kubernetes, or tencentcloud).",
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the account discovery configuration is enabled.",
				Computed:    true,
			},
			"config": schema.StringAttribute{
				Description: "JSON-encoded configuration specific to the provider.",
				Computed:    true,
				Sensitive:   true,
			},
			"schedule": schema.StringAttribute{
				Description: "JSON-encoded schedule information for when the discovery runs.",
				Computed:    true,
			},
			"validity": schema.StringAttribute{
				Description: "JSON-encoded information about the most recent credential validation attempt.",
				Computed:    true,
			},
		},
	}
}

func (d *accountDiscoveryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDiscoveryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountDiscoveryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		AccountDiscovery struct {
			ID          string
			Name        string
			Description string
			Provider    string
			Config      struct {
				TypeName string `graphql:"__typename"`
				// AWS-specific fields
				OrgReadRole   string `graphql:"... on AWSAccountDiscoveryConfig"`
				MemberRole    string `graphql:"... on AWSAccountDiscoveryConfig"`
				CustodianRole string `graphql:"... on AWSAccountDiscoveryConfig"`
				// Azure-specific fields
				TenantID     string `graphql:"... on AzureAccountDiscoveryConfig"`
				ClientID     string `graphql:"... on AzureAccountDiscoveryConfig"`
				ClientSecret string `graphql:"... on AzureAccountDiscoveryConfig"`
				// GCP-specific fields
				OrgID            string   `graphql:"... on GCPAccountDiscoveryConfig"`
				RootFolderIDs    []string `graphql:"... on GCPAccountDiscoveryConfig"`
				ExcludeFolderIDs []string `graphql:"... on GCPAccountDiscoveryConfig"`
				CredentialJSON   string   `graphql:"... on GCPAccountDiscoveryConfig"`
			}
			Schedule struct {
				Suspended bool
			}
			Validity struct {
				Valid   bool
				Message string
			}
		} `graphql:"accountDiscovery(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(data.Name.ValueString()),
	}

	err := d.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account discovery, got error: %s", err))
		return
	}

	if query.AccountDiscovery.Name == "" {
		resp.Diagnostics.AddError("Not Found", "No account discovery found with the specified name")
		return
	}

	// Convert config to JSON based on provider type
	var configJSON string
	switch query.AccountDiscovery.Config.TypeName {
	case "AWSAccountDiscoveryConfig":
		configJSON = fmt.Sprintf(`{
			"org_read_role": %q,
			"member_role": %q,
			"custodian_role": %q
		}`, query.AccountDiscovery.Config.OrgReadRole,
			query.AccountDiscovery.Config.MemberRole,
			query.AccountDiscovery.Config.CustodianRole)
	case "AzureAccountDiscoveryConfig":
		configJSON = fmt.Sprintf(`{
			"tenant_id": %q,
			"client_id": %q,
			"client_secret": %q
		}`, query.AccountDiscovery.Config.TenantID,
			query.AccountDiscovery.Config.ClientID,
			query.AccountDiscovery.Config.ClientSecret)
	case "GCPAccountDiscoveryConfig":
		configJSON = fmt.Sprintf(`{
			"org_id": %q,
			"root_folder_ids": %q,
			"exclude_folder_ids": %q,
			"credential_json": %q
		}`, query.AccountDiscovery.Config.OrgID,
			query.AccountDiscovery.Config.RootFolderIDs,
			query.AccountDiscovery.Config.ExcludeFolderIDs,
			query.AccountDiscovery.Config.CredentialJSON)
	}

	data.ID = types.StringValue(query.AccountDiscovery.ID)
	data.Name = types.StringValue(query.AccountDiscovery.Name)
	data.Description = types.StringValue(query.AccountDiscovery.Description)
	data.CloudProvider = types.StringValue(query.AccountDiscovery.Provider)
	data.Config = types.StringValue(configJSON)
	data.Schedule = types.StringValue(fmt.Sprintf(`{"suspended": %t}`, query.AccountDiscovery.Schedule.Suspended))
	data.Validity = types.StringValue(fmt.Sprintf(`{"valid": %t, "message": %q}`,
		query.AccountDiscovery.Validity.Valid,
		query.AccountDiscovery.Validity.Message))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
