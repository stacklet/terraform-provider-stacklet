# stacklet_binding (Data Source)

Retrieves information about a binding in Stacklet. A binding connects an account group to a policy collection, defining how and when policies should be applied to the accounts in the group.

## Example Usage

```hcl
# Fetch a binding by UUID
data "stacklet_binding" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch a binding by name
data "stacklet_binding" "example" {
  name = "production-security-policies"
}

# Output binding details
output "binding_name" {
  value = data.stacklet_binding.example.name
}

output "binding_percentage_deployed" {
  value = data.stacklet_binding.example.percentage_deployed
}
```

## Argument Reference

* `uuid` - (Optional) The UUID of the binding.
* `name` - (Optional) The name of the binding.

At least one of `uuid` or `name` must be specified.

## Attribute Reference

* `id` - The GraphQL Node ID of the binding.
* `uuid` - The UUID of the binding.
* `name` - The name of the binding.
* `description` - A description of the binding.
* `auto_deploy` - Whether the binding automatically deploys when the policy collection changes.
* `system` - Whether this is a system binding (not user editable).
* `schedule` - The schedule for the binding (e.g., 'rate(1 hour)', 'rate(2 hours)', or cron expression).
* `variables` - JSON-encoded dictionary of values used for policy templating.
* `last_deployed` - The timestamp of the last deployment.
* `account_group_uuid` - The UUID of the account group this binding applies to.
* `policy_collection_uuid` - The UUID of the policy collection this binding applies.
* `percentage_deployed` - The percentage of accounts where the binding is deployed (0-100).
* `is_stale` - Whether the binding has pending changes that need to be deployed. 