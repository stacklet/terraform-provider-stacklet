package datasources

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
)

var (
	_ datasource.DataSource = &bindingDataSource{}
)

func NewBindingDataSource() datasource.DataSource {
	return &bindingDataSource{}
}

type bindingDataSource struct {
	client *graphql.Client
}

type bindingDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	System               types.Bool   `tfsdk:"system"`
	Schedule             types.String `tfsdk:"schedule"`
	Variables            types.String `tfsdk:"variables"`
	LastDeployed         types.String `tfsdk:"last_deployed"`
	AccountGroupUUID     types.String `tfsdk:"account_group_uuid"`
	PolicyCollectionUUID types.String `tfsdk:"policy_collection_uuid"`
	PercentageDeployed   types.Number `tfsdk:"percentage_deployed"`
	IsStale              types.Bool   `tfsdk:"is_stale"`
}

func (d *bindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding"
}

func (d *bindingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a binding by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the binding.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the binding.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the binding.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the binding.",
				Computed:    true,
			},
			"auto_deploy": schema.BoolAttribute{
				Description: "Whether the binding automatically deploys when the policy collection changes.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system binding (not user editable).",
				Computed:    true,
			},
			"schedule": schema.StringAttribute{
				Description: "The schedule for the binding (e.g., 'rate(1 hour)', 'rate(2 hours)', or cron expression).",
				Computed:    true,
			},
			"variables": schema.StringAttribute{
				Description: "JSON-encoded dictionary of values used for policy templating.",
				Computed:    true,
			},
			"last_deployed": schema.StringAttribute{
				Description: "The timestamp of the last deployment.",
				Computed:    true,
			},
			"account_group_uuid": schema.StringAttribute{
				Description: "The UUID of the account group this binding applies to.",
				Computed:    true,
			},
			"policy_collection_uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection this binding applies.",
				Computed:    true,
			},
			"percentage_deployed": schema.NumberAttribute{
				Description: "The percentage of accounts where the binding is deployed (0-100).",
				Computed:    true,
			},
			"is_stale": schema.BoolAttribute{
				Description: "Whether the binding has pending changes that need to be deployed.",
				Computed:    true,
			},
		},
	}
}

func (d *bindingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *bindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bindingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GraphQL query
	var query struct {
		Binding struct {
			ID           string
			UUID         string
			Name         string
			Description  string
			AutoDeploy   bool
			System       bool
			Schedule     string
			Variables    string
			LastDeployed string
			AccountGroup struct {
				UUID string
			}
			PolicyCollection struct {
				UUID string
			}
			PercentageDeployed float64
			IsStale            bool
		} `graphql:"binding(uuid: $uuid, name: $name)"`
	}

	variables := map[string]any{
		"uuid": (*string)(nil),
		"name": (*string)(nil),
	}

	if !data.UUID.IsNull() {
		variables["uuid"] = graphql.String(data.UUID.ValueString())
	}

	if !data.Name.IsNull() {
		variables["name"] = graphql.String(data.Name.ValueString())
	}

	err := d.client.Query(ctx, &query, variables)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read binding, got error: %s", err))
		return
	}

	if query.Binding.UUID == "" {
		resp.Diagnostics.AddError("Not Found", "No binding found with the specified UUID or name")
		return
	}

	data.ID = types.StringValue(query.Binding.ID)
	data.UUID = types.StringValue(query.Binding.UUID)
	data.Name = types.StringValue(query.Binding.Name)
	data.Description = types.StringValue(query.Binding.Description)
	data.AutoDeploy = types.BoolValue(query.Binding.AutoDeploy)
	data.System = types.BoolValue(query.Binding.System)
	data.Schedule = types.StringValue(query.Binding.Schedule)
	data.Variables = types.StringValue(query.Binding.Variables)
	data.LastDeployed = types.StringValue(query.Binding.LastDeployed)
	data.AccountGroupUUID = types.StringValue(query.Binding.AccountGroup.UUID)
	data.PolicyCollectionUUID = types.StringValue(query.Binding.PolicyCollection.UUID)
	data.PercentageDeployed = types.NumberValue(big.NewFloat(query.Binding.PercentageDeployed))
	data.IsStale = types.BoolValue(query.Binding.IsStale)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
