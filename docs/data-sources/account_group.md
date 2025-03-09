# stacklet_account_group (Data Source)

Retrieves information about an account group in Stacklet. Account groups are used to organize and manage collections of cloud accounts.

## Example Usage

```hcl
# Fetch an account group by UUID
data "stacklet_account_group" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch an account group by name
data "stacklet_account_group" "example" {
  name = "production-accounts"
}

# Output account group details
output "group_cloud_provider" {
  value = data.stacklet_account_group.example.cloud_provider
}

output "group_description" {
  value = data.stacklet_account_group.example.description
}
```

## Argument Reference

* `uuid` - (Optional) The UUID of the account group.
* `name` - (Optional) The name of the account group.

At least one of `uuid` or `name` must be specified.

## Attribute Reference

* `id` - The GraphQL Node ID of the account group.
* `uuid` - The UUID of the account group.
* `name` - The name of the account group.
* `description` - The description of the account group.
* `cloud_provider` - The cloud provider for the account group (aws, azure, gcp, kubernetes, or tencentcloud). 