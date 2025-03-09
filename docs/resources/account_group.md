# stacklet_account_group (Resource)

Manages an account group in Stacklet. Account groups are used to organize and manage collections of cloud accounts.

## Example Usage

```hcl
# Create an AWS account group
resource "stacklet_account_group" "production" {
  name           = "production-accounts"
  cloud_provider = "aws"
  description    = "Production AWS accounts"
}

# Create an Azure account group
resource "stacklet_account_group" "development" {
  name           = "development-accounts"
  cloud_provider = "azure"
  description    = "Development Azure accounts"
}
```

## Argument Reference

* `name` - (Required) The name of the account group.
* `cloud_provider` - (Required) The cloud provider for the account group (aws, azure, gcp, kubernetes, or tencentcloud).
* `description` - (Optional) A description of the account group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GraphQL Node ID of the account group.
* `uuid` - The UUID of the account group.

## Import

Account groups can be imported using their UUID:

```shell
terraform import stacklet_account_group.production 00000000-0000-0000-0000-000000000000
``` 