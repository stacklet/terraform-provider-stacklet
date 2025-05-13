# Add an account to a group
resource "stacklet_account_group_mapping" "example" {
  group_uuid     = "00000000-0000-0000-0000-000000000000"
  account_key    = "123456789012"
  cloud_provider = "AWS"
}

# Reference existing account and group
resource "stacklet_account_group_mapping" "example" {
  group_uuid     = data.stacklet_account_group.production.uuid
  account_key    = data.stacklet_account.prod_account.key
  cloud_provider = data.stacklet_account.prod_account.cloud_provider
}

# Add multiple accounts to a group
resource "stacklet_account_group_mapping" "prod_accounts" {
  for_each = toset([
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333",
  ])

  group_uuid   = stacklet_account_group.production.uuid
  account_uuid = each.value
}
