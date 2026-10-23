// Copyright Stacklet, Inc. 2025, 2026

// package api provides access to the GraphQL API.
package api

// API provides access to the GraphQL API.
type API struct {
	Account                 accountAPI
	AccountDiscovery        accountDiscoveryAPI
	AccountGroup            accountGroupAPI
	AccountGroupMapping     accountGroupMappingAPI
	AnthropicIntegration    anthropicIntegrationAPI
	Binding                 bindingAPI
	ConfigurationProfile    configurationProfileAPI
	DatadogIntegration      datadogIntegrationAPI
	GCPIntegration          gcpIntegrationAPI
	Policy                  policyAPI
	PolicyCollection        policyCollectionAPI
	PolicyCollectionMapping policyCollectionMappingAPI
	ReportGroup             reportGroupAPI
	Repository              repositoryAPI
	Role                    roleAPI
	RoleAssignment          roleAssignmentAPI
	SAMLProvider            samlProviderAPI
	SSOGroup                ssoGroupAPI
	System                  systemAPI
	Template                templateAPI
	User                    userAPI
	UserGroup               userGroupAPI
}

// New creates an API wrapper.
func New(config ClientConfig) *API {
	c := newClient(config)
	return &API{
		Account:                 accountAPI{c},
		AccountDiscovery:        accountDiscoveryAPI{c},
		AccountGroup:            accountGroupAPI{c},
		AccountGroupMapping:     accountGroupMappingAPI{c},
		AnthropicIntegration:    anthropicIntegrationAPI{c},
		Binding:                 bindingAPI{c},
		ConfigurationProfile:    configurationProfileAPI{c},
		DatadogIntegration:      datadogIntegrationAPI{c},
		GCPIntegration:          gcpIntegrationAPI{c},
		Policy:                  policyAPI{c},
		PolicyCollection:        policyCollectionAPI{c},
		PolicyCollectionMapping: policyCollectionMappingAPI{c},
		ReportGroup:             reportGroupAPI{c},
		Repository:              repositoryAPI{c},
		Role:                    roleAPI{c},
		RoleAssignment:          roleAssignmentAPI{c},
		SAMLProvider:            samlProviderAPI{c},
		SSOGroup:                ssoGroupAPI{c},
		System:                  systemAPI{c},
		Template:                templateAPI{c},
		User:                    userAPI{c},
		UserGroup:               userGroupAPI{c},
	}
}
