package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &accountDataSource{}
)

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSource struct {
	api *api.API
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
				Description: "The display name for the account.",
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
			"variables": schema.StringAttribute{
				Description: "JSON encoded dict of values used for policy templating.",
				Computed:    true,
			},
		},
	}
}

func (d *accountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.AccountDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := d.api.Account.Read(ctx, data.CloudProvider.ValueString(), data.Key.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(account.ID)
	data.Key = types.StringValue(account.Key)
	data.Name = types.StringValue(account.Name)
	data.ShortName = tftypes.NullableString(account.ShortName)
	data.Description = tftypes.NullableString(account.Description)
	data.CloudProvider = types.StringValue(string(account.Provider))
	data.Path = tftypes.NullableString(account.Path)
	data.Email = tftypes.NullableString(account.Email)
	data.SecurityContext = tftypes.NullableString(account.SecurityContext)
	variablesString, err := tftypes.JSONString(account.Variables)
	if err != nil {
		resp.Diagnostics.AddError("Invalid content for variables", err.Error())
		return
	}
	data.Variables = variablesString
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
