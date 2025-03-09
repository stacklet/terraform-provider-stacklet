# stacklet_sso_group (Resource)

Manages an SSO group configuration in Stacklet. SSO groups allow you to map external SSO provider groups to Stacklet roles and account group access, enabling automated access control management for your users.

## Example Usage

```hcl
# Create an SSO group for administrators
resource "stacklet_sso_group" "admins" {
  name = "Administrators"
  roles = [
    "admin",
    "viewer"
  ]
  account_group_uuids = [
    "00000000-0000-0000-0000-000000000001",  # Production accounts
    "00000000-0000-0000-0000-000000000002"   # Development accounts
  ]
}

# Create an SSO group for developers
resource "stacklet_sso_group" "developers" {
  name = "Developers"
  roles = [
    "viewer"
  ]
  account_group_uuids = [
    "00000000-0000-0000-0000-000000000002"   # Development accounts only
  ]
}

# Output the group IDs
output "admin_group_id" {
  value = stacklet_sso_group.admins.id
}

output "developer_group_id" {
  value = stacklet_sso_group.developers.id
}
```

## Argument Reference

* `name` - (Required) The name that identifies the group in the external SSO provider.
* `roles` - (Required) List of Stacklet API roles automatically granted to SSO users in this group.
* `account_group_uuids` - (Required) List of account group UUIDs whose resources are visible to SSO users in this group.

## Attribute Reference

* `id` - Unique identifier for this SSO group configuration.

## Import

SSO groups can be imported using the group name:

```shell
terraform import stacklet_sso_group.admins Administrators
``` 