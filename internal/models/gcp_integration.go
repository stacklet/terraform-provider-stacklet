// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// GCPIntegrationDataSource is the model for the GCP integration data source.
type GCPIntegrationDataSource struct {
	ID             types.String `tfsdk:"id"`
	Key            types.String `tfsdk:"key"`
	CustomerConfig types.Object `tfsdk:"customer_config"`
	AccessConfig   types.Object `tfsdk:"access_config"`
}

func (m *GCPIntegrationDataSource) Update(integration *api.GCPIntegration) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(string(integration.ID))
	m.Key = types.StringValue(integration.Key)

	customerConfig, d := m.buildCustomerConfig(integration.CustomerConfig)
	errors.AddAttributeDiags(&diags, d, "customer_config")
	m.CustomerConfig = customerConfig

	accessConfig, d := m.buildAccessConfig(integration.AccessConfig)
	errors.AddAttributeDiags(&diags, d, "access_config")
	m.AccessConfig = accessConfig

	return diags
}

func (m GCPIntegrationDataSource) buildCustomerConfig(config api.GCPIntegrationCustomerConfig) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(GCPIntegrationCustomerConfigModel{}.AttributeTypes())

	infrastructure, d := m.buildCustomerInfra(config.Infrastructure)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	organizations, d := typehelpers.ObjectList[GCPIntegrationCustomerOrgModel](
		config.Organizations,
		func(org api.GCPIntegrationCustomerOrganization) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"org_id":      types.StringValue(org.OrgID),
				"folder_ids":  typehelpers.StringsList(org.FolderIDs),
				"project_ids": typehelpers.StringsList(org.ProjectIDs),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	costSources, d := typehelpers.ObjectList[GCPIntegrationCustomerCostSourceModel](
		config.CostSources,
		func(cs api.GCPIntegrationCustomerCostSource) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"billing_table": types.StringValue(cs.BillingTable),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	securityContexts, d := typehelpers.ObjectList[GCPIntegrationCustomerSecurityContextModel](
		config.SecurityContexts,
		func(sc api.GCPIntegrationCustomerSecurityContext) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name":        types.StringValue(sc.Name),
				"extra_roles": typehelpers.StringsList(sc.ExtraRoles),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	terraformModule, d := m.buildTerraformModule(config.TerraformModule)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	return types.ObjectValue(
		GCPIntegrationCustomerConfigModel{}.AttributeTypes(),
		map[string]attr.Value{
			"infrastructure":    infrastructure,
			"organizations":     organizations,
			"cost_sources":      costSources,
			"security_contexts": securityContexts,
			"terraform_module":  terraformModule,
		},
	)
}

func (m GCPIntegrationDataSource) buildCustomerInfra(infra *api.GCPIntegrationCustomerInfrastructure) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(GCPIntegrationCustomerInfraModel{}.AttributeTypes())

	if infra == nil {
		return nullObj, diags
	}

	createProject, d := m.buildCustomerCreateProject(infra.CreateProject)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	return types.ObjectValue(
		GCPIntegrationCustomerInfraModel{}.AttributeTypes(),
		map[string]attr.Value{
			"project_id":        types.StringValue(infra.ProjectID),
			"resource_location": types.StringValue(infra.ResourceLocation),
			"resource_prefix":   types.StringValue(infra.ResourcePrefix),
			"create_project":    createProject,
		},
	)
}

func (m GCPIntegrationDataSource) buildCustomerCreateProject(cp *api.GCPIntegrationCustomerCreateProject) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(GCPIntegrationCustomerCreateProjectModel{}.AttributeTypes())

	if cp == nil {
		return nullObj, diags
	}

	labels, d := typehelpers.ObjectList[Tag](
		cp.Labels,
		func(t api.Tag) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"key":   types.StringValue(t.Key),
				"value": types.StringValue(t.Value),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	return types.ObjectValue(
		GCPIntegrationCustomerCreateProjectModel{}.AttributeTypes(),
		map[string]attr.Value{
			"billing_account_id": types.StringValue(cp.BillingAccountID),
			"org_id":             types.StringPointerValue(cp.OrgID),
			"folder_id":          types.StringPointerValue(cp.FolderID),
			"labels":             labels,
		},
	)
}

func (m GCPIntegrationDataSource) buildTerraformModule(tm *api.TerraformModule) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(TerraformModule{}.AttributeTypes())

	if tm == nil {
		return nullObj, diags
	}

	return types.ObjectValue(
		TerraformModule{}.AttributeTypes(),
		map[string]attr.Value{
			"repository_url": types.StringValue(tm.RepositoryURL),
			"source":         types.StringValue(tm.Source),
			"version":        types.StringPointerValue(tm.Version),
			"variables_json": types.StringValue(tm.VariablesJSON),
		},
	)
}

func (m GCPIntegrationDataSource) buildAccessConfig(config *api.GCPIntegrationAccessConfig) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(GCPIntegrationAccessConfigModel{}.AttributeTypes())

	if config == nil {
		return nullObj, diags
	}

	infrastructure, d := m.buildAccessInfra(config.Infrastructure)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	organizations, d := typehelpers.ObjectList[GCPIntegrationAccessOrgModel](
		config.Organizations,
		func(org api.GCPIntegrationAccessOrganization) (map[string]attr.Value, diag.Diagnostics) {
			folders, d := typehelpers.ObjectList[GCPIntegrationAccessOrgFolderModel](
				org.Folders,
				func(f api.GCPIntegrationAccessOrganizationFolder) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"id":   types.StringValue(f.ID),
						"name": types.StringValue(f.Name),
					}, nil
				},
			)
			if d.HasError() {
				return make(map[string]attr.Value), d
			}

			projects, d := typehelpers.ObjectList[GCPIntegrationAccessOrgProjectModel](
				org.Projects,
				func(p api.GCPIntegrationAccessOrganizationProject) (map[string]attr.Value, diag.Diagnostics) {
					return map[string]attr.Value{
						"id":     types.StringValue(p.ID),
						"number": types.StringValue(p.Number),
					}, nil
				},
			)
			if d.HasError() {
				return make(map[string]attr.Value), d
			}

			return map[string]attr.Value{
				"id":       types.StringValue(org.ID),
				"name":     types.StringValue(org.Name),
				"folders":  folders,
				"projects": projects,
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	costSources, d := typehelpers.ObjectList[GCPIntegrationAccessCostSourceModel](
		config.CostSources,
		func(cs api.GCPIntegrationAccessCostSource) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"billing_table": types.StringValue(cs.BillingTable),
				"location":      types.StringValue(cs.Location),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	securityContexts, d := typehelpers.ObjectList[GCPIntegrationAccessSecurityContextModel](
		config.SecurityContexts,
		func(sc api.GCPIntegrationAccessSecurityContext) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name":        types.StringValue(sc.Name),
				"extra_roles": typehelpers.StringsList(sc.ExtraRoles),
				"principal":   types.StringValue(sc.Principal),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	return types.ObjectValue(
		GCPIntegrationAccessConfigModel{}.AttributeTypes(),
		map[string]attr.Value{
			"infrastructure":    infrastructure,
			"organizations":     organizations,
			"cost_sources":      costSources,
			"security_contexts": securityContexts,
			"roundtrip_digest":  types.StringValue(config.RoundtripDigest),
		},
	)
}

func (m GCPIntegrationDataSource) buildAccessInfra(infra api.GCPIntegrationAccessInfrastructure) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullObj := types.ObjectNull(GCPIntegrationAccessInfraModel{}.AttributeTypes())

	relay, d := types.ObjectValue(
		GCPIntegrationAccessRelayModel{}.AttributeTypes(),
		map[string]attr.Value{
			"oauth_id": types.StringValue(infra.Relay.OAuthID),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	principals, d := types.ObjectValue(
		GCPIntegrationAccessWIFPrincipalsModel{}.AttributeTypes(),
		map[string]attr.Value{
			"read_only":  types.StringValue(infra.WIF.Principals.ReadOnly),
			"cost_query": types.StringValue(infra.WIF.Principals.CostQuery),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	wif, d := types.ObjectValue(
		GCPIntegrationAccessWIFModel{}.AttributeTypes(),
		map[string]attr.Value{
			"audience":   types.StringValue(infra.WIF.Audience),
			"principals": principals,
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj, diags
	}

	return types.ObjectValue(
		GCPIntegrationAccessInfraModel{}.AttributeTypes(),
		map[string]attr.Value{
			"project_id":     types.StringValue(infra.ProjectID),
			"relay":          relay,
			"wif":            wif,
			"baseline_roles": typehelpers.StringsList(infra.BaselineRoles),
		},
	)
}

// GCPIntegrationCustomerConfigModel is the model for GCP integration customer configuration.
type GCPIntegrationCustomerConfigModel struct {
	Infrastructure   types.Object `tfsdk:"infrastructure"`
	Organizations    types.List   `tfsdk:"organizations"`
	CostSources      types.List   `tfsdk:"cost_sources"`
	SecurityContexts types.List   `tfsdk:"security_contexts"`
	TerraformModule  types.Object `tfsdk:"terraform_module"`
}

func (m GCPIntegrationCustomerConfigModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"infrastructure":    types.ObjectType{AttrTypes: GCPIntegrationCustomerInfraModel{}.AttributeTypes()},
		"organizations":     types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationCustomerOrgModel{}.AttributeTypes()}},
		"cost_sources":      types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationCustomerCostSourceModel{}.AttributeTypes()}},
		"security_contexts": types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationCustomerSecurityContextModel{}.AttributeTypes()}},
		"terraform_module":  types.ObjectType{AttrTypes: TerraformModule{}.AttributeTypes()},
	}
}

// GCPIntegrationCustomerInfraModel is the model for GCP integration customer infrastructure.
type GCPIntegrationCustomerInfraModel struct {
	ProjectID        types.String `tfsdk:"project_id"`
	ResourceLocation types.String `tfsdk:"resource_location"`
	ResourcePrefix   types.String `tfsdk:"resource_prefix"`
	CreateProject    types.Object `tfsdk:"create_project"`
}

func (m GCPIntegrationCustomerInfraModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id":        types.StringType,
		"resource_location": types.StringType,
		"resource_prefix":   types.StringType,
		"create_project":    types.ObjectType{AttrTypes: GCPIntegrationCustomerCreateProjectModel{}.AttributeTypes()},
	}
}

// GCPIntegrationCustomerCreateProjectModel is the model for GCP infrastructure project creation configuration.
type GCPIntegrationCustomerCreateProjectModel struct {
	BillingAccountID types.String `tfsdk:"billing_account_id"`
	OrgID            types.String `tfsdk:"org_id"`
	FolderID         types.String `tfsdk:"folder_id"`
	Labels           types.List   `tfsdk:"labels"`
}

func (m GCPIntegrationCustomerCreateProjectModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_account_id": types.StringType,
		"org_id":             types.StringType,
		"folder_id":          types.StringType,
		"labels":             types.ListType{ElemType: types.ObjectType{AttrTypes: Tag{}.AttributeTypes()}},
	}
}

// GCPIntegrationCustomerOrgModel is the model for a GCP customer organization.
type GCPIntegrationCustomerOrgModel struct {
	OrgID      types.String `tfsdk:"org_id"`
	FolderIDs  types.List   `tfsdk:"folder_ids"`
	ProjectIDs types.List   `tfsdk:"project_ids"`
}

func (m GCPIntegrationCustomerOrgModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"org_id":      types.StringType,
		"folder_ids":  types.ListType{ElemType: types.StringType},
		"project_ids": types.ListType{ElemType: types.StringType},
	}
}

// GCPIntegrationCustomerCostSourceModel is the model for a GCP customer cost source.
type GCPIntegrationCustomerCostSourceModel struct {
	BillingTable types.String `tfsdk:"billing_table"`
}

func (m GCPIntegrationCustomerCostSourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_table": types.StringType,
	}
}

// GCPIntegrationCustomerSecurityContextModel is the model for a GCP customer security context.
type GCPIntegrationCustomerSecurityContextModel struct {
	Name       types.String `tfsdk:"name"`
	ExtraRoles types.List   `tfsdk:"extra_roles"`
}

func (m GCPIntegrationCustomerSecurityContextModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"extra_roles": types.ListType{ElemType: types.StringType},
	}
}

// GCPIntegrationAccessConfigModel is the model for GCP integration access configuration.
type GCPIntegrationAccessConfigModel struct {
	Infrastructure   types.Object `tfsdk:"infrastructure"`
	Organizations    types.List   `tfsdk:"organizations"`
	CostSources      types.List   `tfsdk:"cost_sources"`
	SecurityContexts types.List   `tfsdk:"security_contexts"`
	RoundtripDigest  types.String `tfsdk:"roundtrip_digest"`
}

func (m GCPIntegrationAccessConfigModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"infrastructure":    types.ObjectType{AttrTypes: GCPIntegrationAccessInfraModel{}.AttributeTypes()},
		"organizations":     types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationAccessOrgModel{}.AttributeTypes()}},
		"cost_sources":      types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationAccessCostSourceModel{}.AttributeTypes()}},
		"security_contexts": types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationAccessSecurityContextModel{}.AttributeTypes()}},
		"roundtrip_digest":  types.StringType,
	}
}

// GCPIntegrationAccessInfraModel is the model for GCP integration access infrastructure.
type GCPIntegrationAccessInfraModel struct {
	ProjectID     types.String `tfsdk:"project_id"`
	Relay         types.Object `tfsdk:"relay"`
	WIF           types.Object `tfsdk:"wif"`
	BaselineRoles types.List   `tfsdk:"baseline_roles"`
}

func (m GCPIntegrationAccessInfraModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id":     types.StringType,
		"relay":          types.ObjectType{AttrTypes: GCPIntegrationAccessRelayModel{}.AttributeTypes()},
		"wif":            types.ObjectType{AttrTypes: GCPIntegrationAccessWIFModel{}.AttributeTypes()},
		"baseline_roles": types.ListType{ElemType: types.StringType},
	}
}

// GCPIntegrationAccessRelayModel is the model for GCP integration relay configuration.
type GCPIntegrationAccessRelayModel struct {
	OAuthID types.String `tfsdk:"oauth_id"`
}

func (m GCPIntegrationAccessRelayModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"oauth_id": types.StringType,
	}
}

// GCPIntegrationAccessWIFModel is the model for GCP Workload Identity Federation configuration.
type GCPIntegrationAccessWIFModel struct {
	Audience   types.String `tfsdk:"audience"`
	Principals types.Object `tfsdk:"principals"`
}

func (m GCPIntegrationAccessWIFModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"audience":   types.StringType,
		"principals": types.ObjectType{AttrTypes: GCPIntegrationAccessWIFPrincipalsModel{}.AttributeTypes()},
	}
}

// GCPIntegrationAccessWIFPrincipalsModel is the model for WIF service account principals.
type GCPIntegrationAccessWIFPrincipalsModel struct {
	ReadOnly  types.String `tfsdk:"read_only"`
	CostQuery types.String `tfsdk:"cost_query"`
}

func (m GCPIntegrationAccessWIFPrincipalsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"read_only":  types.StringType,
		"cost_query": types.StringType,
	}
}

// GCPIntegrationAccessOrgModel is the model for a GCP accessible organization.
type GCPIntegrationAccessOrgModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Folders  types.List   `tfsdk:"folders"`
	Projects types.List   `tfsdk:"projects"`
}

func (m GCPIntegrationAccessOrgModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"name":     types.StringType,
		"folders":  types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationAccessOrgFolderModel{}.AttributeTypes()}},
		"projects": types.ListType{ElemType: types.ObjectType{AttrTypes: GCPIntegrationAccessOrgProjectModel{}.AttributeTypes()}},
	}
}

// GCPIntegrationAccessOrgFolderModel is the model for a connected organization folder.
type GCPIntegrationAccessOrgFolderModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (m GCPIntegrationAccessOrgFolderModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// GCPIntegrationAccessOrgProjectModel is the model for a connected organization project.
type GCPIntegrationAccessOrgProjectModel struct {
	ID     types.String `tfsdk:"id"`
	Number types.String `tfsdk:"number"`
}

func (m GCPIntegrationAccessOrgProjectModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.StringType,
		"number": types.StringType,
	}
}

// GCPIntegrationAccessCostSourceModel is the model for a GCP access cost source.
type GCPIntegrationAccessCostSourceModel struct {
	BillingTable types.String `tfsdk:"billing_table"`
	Location     types.String `tfsdk:"location"`
}

func (m GCPIntegrationAccessCostSourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_table": types.StringType,
		"location":      types.StringType,
	}
}

// GCPIntegrationAccessSecurityContextModel is the model for a GCP access security context.
type GCPIntegrationAccessSecurityContextModel struct {
	Name       types.String `tfsdk:"name"`
	ExtraRoles types.List   `tfsdk:"extra_roles"`
	Principal  types.String `tfsdk:"principal"`
}

func (m GCPIntegrationAccessSecurityContextModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"extra_roles": types.ListType{ElemType: types.StringType},
		"principal":   types.StringType,
	}
}
