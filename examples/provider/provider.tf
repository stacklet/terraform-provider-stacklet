terraform {
  required_providers {
    stacklet = {
      source = "registry.terraform.io/stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.stacklet.io" # Optional: can also use STACKLET_ENDPOINT env var
  api_key  = "your-api-key"            # Optional: can also use STACKLET_API_KEY env var
}
