// Copyright (c) 2025 - Stacklet, Inc.

package api

// CloudProvider represents a cloud service provider in Stacklet.
type CloudProvider StringEnum

const (
	CloudProviderAWS          = CloudProvider("AWS")
	CloudProviderAzure        = CloudProvider("Azure")
	CloudProviderGCP          = CloudProvider("GCP")
	CloudProviderKubernetes   = CloudProvider("Kubernetes")
	CloudProviderTencentCloud = CloudProvider("TencentCloud")
)

// CLOUD_PROVIDERS is the list of all supported cloud providers.
var CLOUD_PROVIDERS = []CloudProvider{
	CloudProviderAWS,
	CloudProviderAzure,
	CloudProviderGCP,
	CloudProviderKubernetes,
	CloudProviderTencentCloud,
}

// ReportSource represents a report group source.
type ReportSource StringEnum

const (
	ReportSourceBinding = ReportSource("BINDING")
	ReportSourceControl = ReportSource("CONTROL")
	ReportSourcePolicy  = ReportSource("POLICY")
)

// REPORT_SOURCES is the list of all supported report group sources.
var REPORT_SOURCES = []ReportSource{
	ReportSourceBinding,
	ReportSourceControl,
	ReportSourcePolicy,
}
