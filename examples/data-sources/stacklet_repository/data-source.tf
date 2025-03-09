# Fetch a repository by URL
data "stacklet_repository" "by_url" {
  url = "https://github.com/example/repo"
}

# Fetch a repository by name
data "stacklet_repository" "by_name" {
  name = "example-repo"
}

# Fetch a repository by UUID
data "stacklet_repository" "by_uuid" {
  uuid = "12345678-90ab-cdef-1234-567890abcdef"
}

# Example of using the data source attributes
output "repository_info" {
  value = {
    name              = data.stacklet_repository.by_url.name
    description       = data.stacklet_repository.by_url.description
    policy_suffixes   = data.stacklet_repository.by_url.policy_file_suffix
    policy_dirs       = data.stacklet_repository.by_url.policy_directories
    branch            = data.stacklet_repository.by_url.branch_name
    last_scanned      = data.stacklet_repository.by_url.last_scanned
    head_commit       = data.stacklet_repository.by_url.head
    vcs_provider      = data.stacklet_repository.by_url.vcs_provider
    is_system_repo    = data.stacklet_repository.by_url.system
  }
} 