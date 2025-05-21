# Fetch an account group by UUID
data "stacklet_account_group" "by_uuid" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch an account group by name
data "stacklet_account_group" "by_name" {
  name = "production-accounts"
}
