resource "stacklet_configuration_profile_servicenow" "example" {
  endpoint = "https://example.com/servicenow"

  username            = "user"
  password_wo         = "sekret"
  password_wo_version = "1"

  issue_type   = "issue"
  closed_state = "closed"
}
