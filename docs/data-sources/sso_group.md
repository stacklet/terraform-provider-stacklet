# stacklet_sso_group (Data Source)

Retrieves information about an SSO group configuration in Stacklet. SSO groups allow you to manage access control by mapping external SSO provider groups to Stacklet roles and account group access.

## Example Usage

```hcl
# Fetch an SSO group configuration
data "stacklet_sso_group" "example" {
  name = "Administrators"
}

# Output the roles assigned to the group
output "admin_roles" {
  value = data.stacklet_sso_group.example.roles
}

# Output the account groups accessible to the group
output "admin_account_groups" {
  value = data.stacklet_sso_group.example.account_group_uuids
}
```

## Argument Reference

* `name` - (Required) The name that identifies the group in the external SSO provider.

## Attribute Reference

* `id` - Unique identifier for this SSO group configuration.
* `roles` - List of Stacklet API roles automatically granted to SSO users in this group.
* `account_group_uuids` - List of account group UUIDs whose resources are visible to SSO users in this group. 