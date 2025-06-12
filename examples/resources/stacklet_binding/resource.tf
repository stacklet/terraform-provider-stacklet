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

  dry_run = true

  resource_limits = {
    max_count      = 200
    max_percentage = 20.0
    requires_both  = true
  }

  # map keys are unqualified policy names
  policy_resource_limits = {
    policy1 = {
      max_count = 10
    }
    policy2 = {
      max_percentage = 30.0
    }
    policy3 = {
      max_count      = 20
      max_percentage = 50.0
      requires_both  = true
    }
  }

  variables = jsonencode({
    environment = "development"
    severity    = "medium"
  })
}
