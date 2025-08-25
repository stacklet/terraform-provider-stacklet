# Terraform Provider for Stacklet

This Terraform Provider allows you to interact with Stacklet's GraphQL API to
manage your resources through Infrastructure as Code.


## Using the provider

The provider is configured as follows:

```terraform
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.<INSTANCE_NAME>.stacklet.io/"
  api_key  = "<API_KEY>"
}
```

### Environment variables

As an alternative, endpoint and key can be defined as environment variables:

```bash
export STACKLET_ENDPOINT="https://api.<INSTANCE_NAME>.stacklet.io/"
export STACKLET_API_KEY="<API_KEY>"
```

### Login via `stacklet-admin` CLI

The provider can also look up authentication details from the
[`stacklet-admin`](https://github.com/stacklet/stacklet-admin) CLI.

After configuring and logging in to the instance via the CLI (`stacklet-admin
login`), the provider will be able to connect to it without needing to specify
credentials in the configuration or via environment variables.

### Example configuration

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
  endpoint = "https://api.<INSTANCE_NAME>.stacklet.io/"
  api_key  = "<API_KEY>"
}

data "stacklet_policy_collection" "example" {
  name = "aws policies for cis-aws"
}

data "stacklet_policy" "one" {
  name = "aws-neptune-cluster-encrypted-rtc"
}

resource "stacklet_policy_collection" "example" {
  name           = "example-collection"
  cloud_provider = "AWS"
  description    = "Example policy collection"
  auto_update    = true
}

resource "stacklet_account_group" "example" {
  name           = "example-account-group"
  cloud_provider = "AWS"
  description    = "Example account group"
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
  name           = "test-acccount"
  short_name     = "tftest"
  description    = "Test account"
  email          = "cloud-team@example.com"
}

resource "stacklet_binding" "binding" {
  name                   = "test-binding"
  description            = "Created with terraform"
  account_group_uuid     = stacklet_account_group.example.uuid
  policy_collection_uuid = stacklet_policy_collection.example.uuid
}

data "stacklet_binding" "binding" {
  name = "AWS Posture"
}
```


## Local development

For local development, make sure you have the tools declared in the
[`.tool-versions`](./.tool-versions) file installed.

### Building the provider

1. Clone the repository:
   ```bash
   git clone https://github.com/stacklet/terraform-provider-stacklet.git
   cd terraform-provider-stacklet
   ```

2. Build the provider:
   ```bash
   just build
   ```

### Running locally built provider

To run the locally built copy of the provider, terraform must be configured as
follows:

1. Override the provider location for development, by creating a
   `~/.terraformrc` with the following content:

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

2. Declare the provider in your terraform configuration as

```terraform
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.<INSTANCE_NAME>.stacklet.io/"  # Or set STACKLET_ENDPOINT env var
  api_key  = "<API_KEY>"                                 # Or set STACKLET_API_KEY env var
}
```

3. Run `terraform plan` or `terraform apply` with the local resources configuration.


**Note**: `terraform init` must not be run when working with a locally installed provider.

### Debugging

Debug messages and output are not visible when running the provider directly
from terraform.  To enable debug:

1. Run `./terraform-provider-stacklet -debug` in one terminal.

2. In a separate terminal, export the value for the `TF_REATTACH_PROVIDERS`
   variable provided in the output of the previous command, and run
   `terraform`.


## Release process

1. Update the [Changelog](./CHANGELOG.md) with an entry for the new release.
2. Create a release tag with `just tag-release X.Y.Z`.
3. Push the tag upstream. This will start the Release workflow which creates
   the release on GitHub and builds packages.
4. Once the workflow has run successfully, publish the release. The Terraform
   registry will pick up the new release automatically.
