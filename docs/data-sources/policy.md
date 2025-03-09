# stacklet_policy (Data Source)

Retrieves information about a policy in Stacklet. Policies define the rules and configurations that are applied to cloud resources for governance, security, and compliance.

## Example Usage

```hcl
# Fetch a policy by UUID
data "stacklet_policy" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch a policy by name
data "stacklet_policy" "example" {
  name = "s3-bucket-encryption"
}

# Output policy details
output "policy_version" {
  value = data.stacklet_policy.example.version
}

output "policy_enabled" {
  value = data.stacklet_policy.example.enabled
}
```

## Argument Reference

* `uuid` - (Optional) The UUID of the policy.
* `name` - (Optional) The name of the policy.

At least one of `uuid` or `name` must be specified.

## Attribute Reference

* `id` - The GraphQL Node ID of the policy.
* `uuid` - The UUID of the policy.
* `name` - The name of the policy.
* `description` - The description of the policy.
* `cloud_provider` - The cloud provider for the policy (aws, azure, gcp, kubernetes, or tencentcloud).
* `version` - The version of the policy.
* `collection_id` - The ID of the policy collection this policy belongs to.
* `enabled` - Whether the policy is enabled. 