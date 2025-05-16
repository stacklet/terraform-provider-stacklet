# AWS account discovery configuration
resource "stacklet_account_discovery" "aws_example" {
  name           = "aws-production"
  description    = "AWS production organization discovery"
  cloud_provider = "AWS"
  org_read_role  = "arn:aws:iam::123456789012:role/StackletOrgReadRole"
  member_role    = "arn:aws:iam::{Id}:role/StackletAssetDBRole"
  custodian_role = "arn:aws:iam::{id}:role/StackletCustodianRole-{tag}"
  suspended      = false
}

# Azure account discovery configuration
resource "stacklet_account_discovery" "azure_example" {
  name           = "azure-production"
  description    = "Azure production tenant discovery"
  cloud_provider = "Azure"
  tenant_id      = "00000000-0000-0000-0000-000000000000"
  client_id      = "11111111-1111-1111-1111-111111111111"
  client_secret  = "your-client-secret"
  suspended      = false
}

# GCP account discovery configuration
resource "stacklet_account_discovery" "gcp_example" {
  name               = "gcp-production"
  description        = "GCP production organization discovery"
  cloud_provider     = "GCP"
  org_id             = "123456789012"
  root_folder_ids    = ["folders/123456", "folders/789012"]
  exclude_folder_ids = ["folders/345678"]
  credential_json = jsonencode({
    "type" : "service_account",
    "project_id" : "your-project",
    "private_key_id" : "key-id",
    "private_key" : "-----BEGIN PRIVATE KEY-----\nYOUR-PRIVATE-KEY\n-----END PRIVATE KEY-----\n",
    "client_email" : "service-account@your-project.iam.gserviceaccount.com",
    "client_id" : "client-id",
    "auth_uri" : "https://accounts.google.com/o/oauth2/auth",
    "token_uri" : "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url" : "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url" : "https://www.googleapis.com/robot/v1/metadata/x509/service-account%40your-project.iam.gserviceaccount.com"
  })
  suspended = false
}
