# Terraform Provider for Stacklet

This Terraform Provider allows you to interact with Stacklet's GraphQL API to manage your resources through Infrastructure as Code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository:
   ```bash
   git clone https://github.com/stacklet/terraform-provider-stacklet.git
   cd terraform-provider-stacklet
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the provider:
   ```bash
   go build -o terraform-provider-stacklet
   ```

4. Install the provider locally for development:
   ```bash
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/stacklet/stacklet/0.1.0/darwin_amd64/
   cp terraform-provider-stacklet ~/.terraform.d/plugins/registry.terraform.io/stacklet/stacklet/0.1.0/darwin_amd64/
   ```

## Using the provider

### Local Development

When developing locally, you can use the provider by configuring Terraform to use your local build:

```hcl
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
      version = "0.1.0"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.dev.stacklet.dev"     # Or use STACKLET_ENDPOINT env var
  api_key  = "your-api-key"                     # Or use STACKLET_API_KEY env var
}
```

### Environment Variables

The provider can be configured using environment variables:

```bash
export STACKLET_ENDPOINT="https://api.stacklet.io/graphql"
export STACKLET_API_KEY="your-api-key"
```

### Testing with Local Provider

1. After building and installing the provider locally (see [Building The Provider](#building-the-provider)), create a test directory:
   ```bash
   mkdir test
   cd test
   ```

2. Create a test configuration file (e.g., `main.tf`):
   ```hcl
   terraform {
     required_providers {
       stacklet = {
         source = "stacklet/stacklet"
         version = "0.1.0"
       }
     }
   }

   # Add your test resources here
   ```

3. Initialize and test:
   ```bash
   terraform init
   terraform plan
   ```
