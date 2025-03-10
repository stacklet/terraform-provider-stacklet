# stacklet_account (Data Source)

Retrieves information about a cloud account in Stacklet. This data source allows you to fetch details about accounts managed by Stacklet across different cloud providers.

## Example Usage

```hcl
# Fetch an AWS account
data "stacklet_account" "aws_prod" {
  cloud_provider = "aws"
  key           = "123456789012"  # AWS account ID
}

# Fetch an Azure subscription
data "stacklet_account" "azure_dev" {
  cloud_provider = "azure"
  key           = "00000000-0000-0000-0000-000000000000"  # Azure subscription ID
}

# Fetch a GCP project
data "stacklet_account" "gcp_staging" {
  cloud_provider = "gcp"
  key           = "my-project-id"  # GCP project ID
}

# Fetch a Tencent Cloud account
data "stacklet_account" "tencent_prod" {
  cloud_provider = "tencentcloud"
  key           = "1234567890"  # Tencent Cloud account ID
}

# Output account details
output "account_name" {
  value = data.stacklet_account.aws_prod.name
}

output "account_variables" {
  value = data.stacklet_account.aws_prod.variables
}
```

## Argument Reference

* `cloud_provider` - (Required) The cloud provider for the account (aws, azure, gcp, or tencentcloud).
* `key` - (Required) The unique identifier for the account within the cloud provider:
  * For AWS: The AWS account ID
  * For Azure: The subscription ID
  * For GCP: The project ID
  * For Tencent Cloud: The account ID

## Attribute Reference

* `id` - The GraphQL Node ID of the account.
* `name` - The display name of the account.
* `short_name` - A shorter display name for the account.
* `description` - A description of the account.
* `path` - The path of the account in the cloud provider's hierarchy.
* `email` - The email address associated with the account.
* `security_context` - JSON-encoded security context for the account.
* `active` - Whether the account is active.
* `variables` - JSON-encoded dictionary of values used for policy templating. 