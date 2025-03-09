# Create an AWS account
resource "stacklet_account" "example" {
  cloud_provider   = "aws"
  key             = "123456789012"
  name            = "Production Account"
  short_name      = "prod"
  description     = "Main production AWS account"
  path            = "/production"
  email           = "cloud-admin@example.com"
  security_context = jsonencode({
    role_arn = "arn:aws:iam::123456789012:role/StackletRole"
  })
  active          = true
  variables       = jsonencode({
    environment = "production"
    team       = "platform"
  })
}

# Output the account ID
output "account_id" {
  value = stacklet_account.example.id
} 