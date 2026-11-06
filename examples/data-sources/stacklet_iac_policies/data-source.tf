# Fetch policies from an IaC policy source
data "stacklet_iac_policies" "from_source" {
  source = "example"
}

# Fetch policies in an IaC policy collection
data "stacklet_iac_policies" "from_collection" {
  collection = "example"
}
