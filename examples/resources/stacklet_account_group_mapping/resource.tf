# Add an existing account to an existing group
resource "stacklet_account_group_mapping" "example" {
  group_uuid     = data.stacklet_account_group.production.uuid
  account_key    = data.stacklet_account.prod_account.key
  cloud_provider = data.stacklet_account.prod_account.cloud_provider
}

data "stacklet_account_group" "production" {
  name = "production-accounts"
}

data "stacklet_account" "prod_account" {
  cloud_provider = "aws"
  key            = "123456789012"
}

locals {
  azure_accounts = [
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333",
  ]
}

resource "stacklet_account_group_mapping" "prod_accounts" {
  for_each = local.azure_accounts

  group_uuid     = stacklet_account_group.production.uuid
  cloud_provider = "azure"
  account_uuid   = each.value
}
