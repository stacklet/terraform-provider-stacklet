# stacklet_policy_collection (Data Source)

Retrieves information about a policy collection in Stacklet. Policy collections are groups of policies that can be applied to cloud resources, organized by cloud provider.

## Example Usage

```hcl
# Fetch a policy collection by UUID
data "stacklet_policy_collection" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch a policy collection by name
data "stacklet_policy_collection" "example" {
  name = "aws-security-policies"
}

# Output policy collection details
output "collection_provider" {
  value = data.stacklet_policy_collection.example.cloud_provider
}

output "collection_auto_update" {
  value = data.stacklet_policy_collection.example.auto_update
}
```

## Argument Reference

* `uuid` - (Optional) The UUID of the policy collection.
* `name` - (Optional) The name of the policy collection.

At least one of `uuid` or `name` must be specified.

## Attribute Reference

* `id` - The GraphQL Node ID of the policy collection.
* `uuid` - The UUID of the policy collection.
* `name` - The name of the policy collection.
* `description` - The description of the policy collection.
* `cloud_provider` - The cloud provider for the policy collection (aws, azure, or gcp).
* `auto_update` - Whether the policy collection automatically updates policy versions.
* `system` - Whether this is a system policy collection (not user editable).
* `repository` - The repository URL if this collection was created from a repo control file. 