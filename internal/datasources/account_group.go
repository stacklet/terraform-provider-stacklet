package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/helpers"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &accountGroupDataSource{}
)

func NewAccountGroupDataSource() datasource.DataSource {
	return &accountGroupDataSource{}
}

type accountGroupDataSource struct {
	api *api.API
}

func (d *accountGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group"
}

func (d *accountGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve an account group by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account group.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the account group, alternative to the name.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the account group, alternative to the UUID.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the account group.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account group.",
				Computed:    true,
			},
			"regions": schema.ListAttribute{
				Description: "The regions for the account group.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *accountGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *accountGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.AccountGroupDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.UUID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of UUID and name must be set")
		return
	}

	account_group, err := d.api.AccountGroup.Read(ctx, data.UUID.ValueString(), data.Name.ValueString())
	if err != nil {
		helpers.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(account_group.ID)
	data.UUID = types.StringValue(account_group.UUID)
	data.Name = types.StringValue(account_group.Name)
	data.Description = tftypes.NullableString(account_group.Description)
	data.CloudProvider = types.StringValue(account_group.Provider)
	data.Regions = tftypes.StringsList(account_group.Regions)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
