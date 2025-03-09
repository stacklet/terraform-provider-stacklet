# Terraform Provider for Stacklet

This Terraform Provider allows you to interact with Stacklet's GraphQL API to manage your resources through Infrastructure as Code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using `go build -o terraform-provider-stacklet`

## Using the provider

```hcl
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.stacklet.io/graphql"  # Or use STACKLET_ENDPOINT env var
  username = "your-username"                     # Or use STACKLET_USERNAME env var
  api_key  = "your-api-key"                     # Or use STACKLET_API_KEY env var
}
```

### Environment Variables

The provider can be configured using environment variables:

```bash
export STACKLET_ENDPOINT="https://api.stacklet.io/graphql"
export STACKLET_USERNAME="your-username"
export STACKLET_API_KEY="your-api-key"
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build -o terraform-provider-stacklet`. This will build the provider and put the provider binary in the current directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. 