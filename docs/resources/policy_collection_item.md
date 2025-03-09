# stacklet_policy_collection_item (Resource)

Manages a policy within a policy collection in Stacklet. This resource allows you to add or remove policies from collections.

## Example Usage

```hcl
# Add a policy to a collection
resource "stacklet_policy_collection_item" "example" {
  collection_uuid = "00000000-0000-0000-0000-000000000000"
  policy_uuid     = "11111111-1111-1111-1111-111111111111"
}

# Reference existing policy and collection
resource "stacklet_policy_collection_item" "example" {
  collection_uuid = data.stacklet_policy_collection.security.uuid
  policy_uuid     = data.stacklet_policy.s3_encryption.uuid
}
```

## Argument Reference

* `collection_uuid` - (Required) The UUID of the policy collection.
* `policy_uuid` - (Required) The UUID of the policy to add to the collection.

## Attribute Reference

* `id` - The ID of the policy collection item, formatted as `{collection_uuid}:{policy_uuid}`. 