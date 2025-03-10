# stacklet_account (Resource)

Manages a cloud account in Stacklet. This resource allows you to create, update, and delete accounts from different cloud providers in Stacklet's platform.

## Example Usage

```hcl
# Configure an AWS account
resource "stacklet_account" "aws_prod" {
  cloud_provider = "aws"
  key           = "123456789012"  # AWS account ID
  name          = "Production AWS Account"
  short_name    = "prod"
  description   = "Main production AWS account"
  email         = "cloud-team@example.com"
  active        = true
  
  variables = jsonencode({
    environment = "production"
    team        = "platform"
    cost_center = "12345"
  })
}

# Configure an Azure subscription
resource "stacklet_account" "azure_dev" {
  cloud_provider = "azure"
  key           = "00000000-0000-0000-0000-000000000000"  # Azure subscription ID
  name          = "Development Azure Subscription"
  short_name    = "dev"
  description   = "Development environment in Azure"
  path          = "/development"
  active        = true
}

# Configure a GCP project
resource "stacklet_account" "gcp_staging" {
  cloud_provider = "gcp"
  key           = "my-project-id"  # GCP project ID
  name          = "Staging GCP Project"
  description   = "Staging environment in GCP"
  
  security_context = jsonencode({
    "org_id": "1234567890"
  })
}

# Configure a Tencent Cloud account
resource "stacklet_account" "tencent_prod" {
  cloud_provider = "tencentcloud"
  key           = "1234567890"  # Tencent Cloud account ID
  name          = "Production Tencent Cloud Account"
  description   = "Production environment in Tencent Cloud"
  path          = "/tencent/production"
}
```

## Argument Reference

* `cloud_provider` - (Required) The cloud provider for the account (aws, azure, gcp, or tencentcloud). This value cannot be changed after creation.
* `key` - (Required) The unique identifier for the account within the cloud provider:
  * For AWS: The AWS account ID
  * For Azure: The subscription ID
  * For GCP: The project ID
  * For Tencent Cloud: The account ID
* `name` - (Required) The display name of the account.
* `short_name` - (Optional) A shorter display name for the account.
* `description` - (Optional) A description of the account.
* `path` - (Optional) The path of the account in the cloud provider's hierarchy.
* `email` - (Optional) The email address associated with the account.
* `security_context` - (Optional) JSON-encoded security context for the account.
* `active` - (Optional) Whether the account is active. Defaults to true.
* `variables` - (Optional) JSON-encoded dictionary of values used for policy templating.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GraphQL Node ID of the account.

## Import

Accounts can be imported using a combination of the cloud provider and key, separated by a colon:

```shell
terraform import stacklet_account.aws_prod aws:123456789012
terraform import stacklet_account.azure_dev azure:00000000-0000-0000-0000-000000000000
terraform import stacklet_account.gcp_staging gcp:my-project-id
terraform import stacklet_account.tencent_prod tencentcloud:1234567890
``` 