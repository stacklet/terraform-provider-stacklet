resource "stacklet_account_discovery_gcp" "example" {
  name               = "gcp"
  description        = "Discovery of GCP accounts"
  org_id             = "1234567890"
  root_folder_ids    = ["folders/123456", "folders/789012"]
  exclude_folder_ids = ["folders/345678"]
  credential_json_wo = jsonencode({
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
  # the version can be any string, as long as it's changed from the previous
  # value whenever the credential_json_wo field needs to be updated
  credential_json_wo_version = "1"
}
