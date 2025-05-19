data "stacklet_account_group" "development" {
  name = "development"
}

data "stacklet_policy_collection" "compliance" {
  name = "compliance"
}

resource "stacklet_binding" "example" {
  name        = "development-compliance"
  description = "Compliance policies for development accounts"

  account_group_uuid     = data.stacklet_account_group.development.uuid
  policy_collection_uuid = data.stacklet_policy_collection.compliance.uuid

  auto_deploy = true
  schedule    = "rate(12 hours)"

  variables = jsonencode({
    environment = "development"
    severity    = "medium"
  })
}
