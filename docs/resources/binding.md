# stacklet_binding (Resource)

Manages a binding in Stacklet. A binding connects an account group to a policy collection, defining how and when policies should be applied to the accounts in the group.

## Example Usage

```hcl
# Create a binding between an account group and a policy collection
resource "stacklet_binding" "example" {
  name                = "production-security-policies"
  description         = "Security policies for production accounts"
  auto_deploy         = true
  schedule            = "rate(1 hour)"
  
  account_group_uuid     = "00000000-0000-0000-0000-000000000000"
  policy_collection_uuid = "11111111-1111-1111-1111-111111111111"
  
  variables = jsonencode({
    environment = "production"
    severity    = "high"
  })
  
  deploy = true  # Deploy immediately after creation
}

# Reference existing account group and policy collection
resource "stacklet_binding" "example" {
  name                = "development-compliance"
  description         = "Compliance policies for development accounts"
  
  account_group_uuid     = data.stacklet_account_group.development.uuid
  policy_collection_uuid = data.stacklet_policy_collection.compliance.uuid
  
  auto_deploy = true
  schedule    = "rate(12 hours)"
}
```

## Argument Reference

* `name` - (Required) The name of the binding.
* `description` - (Optional) A description of the binding.
* `auto_deploy` - (Optional) Whether the binding should automatically deploy when the policy collection changes.
* `schedule` - (Optional) The schedule for the binding (e.g., 'rate(1 hour)', 'rate(2 hours)', or cron expression).
* `variables` - (Optional) JSON-encoded dictionary of values used for policy templating.
* `account_group_uuid` - (Required) The UUID of the account group this binding applies to. This value cannot be changed after creation.
* `policy_collection_uuid` - (Required) The UUID of the policy collection this binding applies. This value cannot be changed after creation.
* `deploy` - (Optional) Whether to deploy the binding immediately after creation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GraphQL Node ID of the binding.
* `uuid` - The UUID of the binding.

## Import

Bindings can be imported using their UUID:

```shell
terraform import stacklet_binding.example 00000000-0000-0000-0000-000000000000
``` 