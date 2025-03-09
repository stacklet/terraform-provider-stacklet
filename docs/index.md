# Stacklet Provider

The Stacklet provider is used to interact with Stacklet's cloud governance platform. The provider allows you to manage various resources such as repositories, accounts, account discoveries, SSO groups, policies, and policy collections.

## Example Usage

```hcl
terraform {
  required_providers {
    stacklet = {
      source = "registry.terraform.io/stacklet/stacklet"
    }
  }
}

provider "stacklet" {
  endpoint = "https://api.stacklet.io"  # Optional: can also use STACKLET_ENDPOINT env var
  api_key  = "your-api-key"            # Optional: can also use STACKLET_API_KEY env var
}
```

## Authentication

The provider supports authentication via an API key. You can provide the API key in two ways:

1. Via the provider configuration:
```hcl
provider "stacklet" {
  api_key = "your-api-key"
}
```

2. Via environment variables:
```sh
export STACKLET_API_KEY="your-api-key"
```

## Provider Configuration

### Required

- `api_key` (String, Sensitive) - The API key for Stacklet authentication. Can also be provided via the `STACKLET_API_KEY` environment variable.

### Optional

- `endpoint` (String) - The endpoint URL of the Stacklet GraphQL API. Can also be provided via the `STACKLET_ENDPOINT` environment variable. Defaults to `https://api.stacklet.io`.

## Resources

- [`stacklet_account`](./resources/account.md) - Manages a cloud account in Stacklet.
- [`stacklet_account_discovery`](./resources/account_discovery.md) - Manages account discovery configurations for cloud providers.
- [`stacklet_policy_collection`](./resources/policy_collection.md) - Manages policy collections for cloud providers.
- [`stacklet_repository`](./resources/repository.md) - Manages a Git repository in Stacklet.
- [`stacklet_sso_group`](./resources/sso_group.md) - Manages SSO group configurations.

## Data Sources

- [`stacklet_account`](./data-sources/account.md) - Retrieves information about a cloud account.
- [`stacklet_account_discovery`](./data-sources/account_discovery.md) - Retrieves information about an account discovery configuration.
- [`stacklet_policy`](./data-sources/policy.md) - Retrieves information about a policy.
- [`stacklet_policy_collection`](./data-sources/policy_collection.md) - Retrieves information about a policy collection.
- [`stacklet_repository`](./data-sources/repository.md) - Retrieves information about a Git repository.
- [`stacklet_sso_group`](./data-sources/sso_group.md) - Retrieves information about an SSO group configuration. 