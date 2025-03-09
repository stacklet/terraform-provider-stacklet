# stacklet_repository (Data Source)

Retrieves information about a Git repository in Stacklet. This data source allows you to fetch details about repositories used for storing and managing policies.

## Example Usage

```hcl
# Fetch a repository by name
data "stacklet_repository" "example" {
  name = "my-policies"
}

# Fetch a repository by URL
data "stacklet_repository" "example" {
  url = "https://github.com/organization/my-policies"
}

# Output repository details
output "repo_uuid" {
  value = data.stacklet_repository.example.uuid
}

output "repo_policy_dirs" {
  value = data.stacklet_repository.example.policy_directories
}
```

## Argument Reference

* `uuid` - (Optional) The UUID of the repository.
* `name` - (Optional) The name of the repository.
* `url` - (Optional) The URL of the repository.

At least one of `uuid`, `name`, or `url` must be specified.

## Attribute Reference

* `uuid` - The UUID of the repository.
* `name` - The name of the repository.
* `url` - The URL of the repository.
* `description` - A description of the repository.
* `policy_file_suffix` - List of file suffixes used for policy scanning.
* `policy_directories` - List of directories that are scanned for policies.
* `branch_name` - The branch used for scanning policies.
* `auth_user` - The user used to access the repository.
* `has_auth_token` - Whether the repository has an auth token configured.
* `has_ssh_private_key` - Whether the repository has an SSH private key configured.
* `has_ssh_passphrase` - Whether the repository has an SSH passphrase configured.
* `head` - The head commit that was processed.
* `last_scanned` - The ISO format datetime of when the repo was last scanned.
* `vcs_provider` - The provider of the repository (e.g., 'github', 'gitlab').
* `system` - Whether this is a system repository (not user editable). 