# Create an API key with description
resource "stacklet_api_key" "example" {
  description = "API key for automated scripts"
}

# Create an API key with expiration
resource "stacklet_api_key" "expiring" {
  description = "Temporary API key for testing"
  expires_at  = "2026-12-31T23:59:59+00:00"
}

# Get the key secret
output "api_key_secret" {
  value     = stacklet_api_key.example.secret
  sensitive = true
}
