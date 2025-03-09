# stacklet_account_group_item (Resource)

Manages an account within an account group in Stacklet. This resource allows you to add or remove accounts from groups.

## Example Usage

```hcl
# Add an account to a group
resource "stacklet_account_group_item" "example" {
  group_uuid   = "00000000-0000-0000-0000-000000000000"
  account_uuid = "11111111-1111-1111-1111-111111111111"
}

# Reference existing account and group
resource "stacklet_account_group_item" "example" {
  group_uuid   = data.stacklet_account_group.production.uuid
  account_uuid = data.stacklet_account.prod_account.uuid
}

# Add multiple accounts to a group
resource "stacklet_account_group_item" "prod_accounts" {
  for_each = toset([
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333",
  ])
  
  group_uuid   = stacklet_account_group.production.uuid
  account_uuid = each.value
}
```

## Argument Reference

* `group_uuid` - (Required) The UUID of the account group.
* `account_uuid` - (Required) The UUID of the account to add to the group.

## Attribute Reference

* `id` - The ID of the account group item, formatted as `{group_uuid}:{account_uuid}`. 