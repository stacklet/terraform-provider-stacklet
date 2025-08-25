resource "stacklet_configuration_profile_jira" "jira" {
  url  = "https://example.com/jira"
  user = "username"

  api_key_wo         = "seKret"
  api_key_wo_version = "1"

  project {
    name          = "PRJ1"
    project       = "123"
    issue_type    = "Issue"
    closed_status = "Closed"
  }

  project {
    name          = "PRJ2"
    project       = "456"
    issue_type    = "Task"
    closed_status = "Done"
  }
}
