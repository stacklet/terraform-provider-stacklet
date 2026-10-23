resource "stacklet_iac_policy_collection" "example" {
  name = "example"
}

resource "stacklet_iac_project" "example" {
  name = "example"
  collections = [
    stacklet_iac_policy_collection.example.name,
  ]
  read_groups = [
    "developers",
  ]
}
