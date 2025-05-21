# Terraform Provider for Stacklet

This Terraform Provider allows you to interact with Stacklet's GraphQL API to manage your resources through Infrastructure as Code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.11
- [Go](https://golang.org/doc/install) >= 1.24.2

## Building The Provider

1. Clone the repository:
   ```bash
   git clone https://github.com/stacklet/terraform-provider-stacklet.git
   cd terraform-provider-stacklet
   ```

2. Build the provider:
   ```bash
   just build
   ```

## Using the provider

### Local Development

When developing locally, you can use the provider by configuring Terraform to use your local build:

1. Override the provider location for development, by creating a `~/.terraformrc` with the following content:

```terraform
provider_installation {
  dev_overrides {
      "stacklet/stacklet" = "<absolute path to the repository directory>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use the
  # dev_overrides block, and so no other providers will be available.
  direct {}
}
```

2. Declare the provider in the terraform file

```terraform
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.<myinstance>.stacklet.io/"  # Or use STACKLET_ENDPOINT env var
  api_key  = "your-api-key"                           # Or use STACKLET_API_KEY env var
}
```

3. Run `terraform plan` or `terraform apply` with the local resources configuration.


**Note**: `terraform init` must not be run when working with a locally installed provider.


#### Debugging the provider

Debug messages and output are not visible when running the provider directly from terraform.
To enable debug

1. run `./terraform-provider-stacklet -debug`

2. export the `TF_REATTACH_PROVIDERS` variable provided in the output in the shell where `terraform` is run


### Environment Variables

The provider can be configured using environment variables:

```bash
export STACKLET_ENDPOINT="https://api.<myinstance>.stacklet.io/"
export STACKLET_API_KEY="your-api-key"
```

### Login via stacklet-admin CLI

The provider can look up authentication details from the [`stacklet-admin`](https://github.com/stacklet/stacklet-admin) CLI.

After configuring and logging in to the instance via the CLI (`stacklet-admin login`), the provider will be able to connect
to it.


### Example Terraform

Below is a full example of a configuration to create a few resources in Stacklet.

```terraform
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.<myinstance>.stacklet.io/"
  api_key  = "$api_key_here"
}

data "stacklet_policy_collection" "example" {
  name = "aws policies for cis-aws"
}

data "stacklet_policy" "one" {
  name = "aws-neptune-cluster-encrypted-rtc"
}

resource "stacklet_policy_collection" "example" {
  name           = "tf-cursor-example-collection"
  cloud_provider = "AWS"
  description    = "Example policy collection"
  auto_update    = true
}

resource "stacklet_account_group" "example" {
  name           = "tf-cusror-example-account-group"
  cloud_provider = "AWS"
  description    = "test account group from terraform"
  regions        = ["us-east-1", "us-east-2"]
}

data "stacklet_account" "one" {
  cloud_provider = "AWS"
  key            = "123456789012"
}

resource "stacklet_account_group_mapping" "one" {
  group_uuid     = stacklet_account_group.example.uuid
  account_key    = data.stacklet_account.one.key
}

resource "stacklet_policy_collection_mapping" "one" {
  collection_uuid = stacklet_policy_collection.example.uuid
  policy_uuid     = data.stacklet_policy.one.uuid
  policy_version  = 2
}

resource "stacklet_account" "two" {
  cloud_provider = "AWS"
  key            = "000000000000" # AWS account ID
  name           = "test-tf-acccount"
  short_name     = "tftest"
  description    = "Test account from terraform testing update"
  email          = "cloud-team@example.com"
}

resource "stacklet_binding" "binding" {
  name                   = "test-tf-cursor-binding"
  description            = "created with terraform update"
  account_group_uuid     = stacklet_account_group.example.uuid
  policy_collection_uuid = stacklet_policy_collection.example.uuid
}

data "stacklet_binding" "binding" {
  name = "AWS Posture"
}
```
