# stacklet_repository (Resource)

Manages a Git repository in Stacklet. This resource allows you to configure repositories that store your policies and control how they are scanned and processed.

## Example Usage

```hcl
# Configure a public repository
resource "stacklet_repository" "public" {
  name        = "public-policies"
  url         = "https://github.com/organization/public-policies"
  description = "Public repository containing our policies"
  
  branch_name = "main"
  policy_directories = [
    "policies/aws",
    "policies/azure"
  ]
  policy_file_suffix = [
    ".yaml",
    ".yml"
  ]
}

# Configure a private repository with authentication
resource "stacklet_repository" "private" {
  name        = "private-policies"
  url         = "https://github.com/organization/private-policies"
  description = "Private repository containing sensitive policies"
  
  auth_user  = "git-user"
  auth_token = "ghp_xxxxxxxxxxxx"  # or use variables: var.github_token
  
  deep_import = true  # Scan repository from the beginning
}

# Configure a repository with SSH authentication
resource "stacklet_repository" "ssh" {
  name        = "ssh-policies"
  url         = "git@github.com:organization/ssh-policies.git"
  description = "Repository accessed via SSH"
  
  ssh_private_key = file("~/.ssh/id_rsa")  # or use variables: var.ssh_key
  ssh_passphrase = "optional-passphrase"   # if the key is encrypted
}
```

## Argument Reference

* `name` - (Required) The name of the repository.
* `url` - (Required) The URL of the repository. This value cannot be changed after creation.
* `description` - (Optional) A description of the repository.
* `policy_file_suffix` - (Optional) Override the default suffix options ['.yaml', '.yml']. This could allow specifying ['.json'] to process other files.
* `policy_directories` - (Optional) If set, only directories that match the list will be scanned for policies.
* `branch_name` - (Optional) If set, use the specified branch name when scanning for policies rather than the repository default.
* `auth_user` - (Optional) The user to use to access the repository if it is private.
* `auth_token` - (Optional, Sensitive) The token for the user to use to access the repository if it is private.
* `ssh_private_key` - (Optional, Sensitive) SSH private key for repository authentication.
* `ssh_passphrase` - (Optional, Sensitive) Passphrase for the SSH private key.
* `deep_import` - (Optional) If true, scan repository from the beginning. If false, only scan the tip.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - The UUID of the repository.

## Import

Repositories can be imported using their URL:

```shell
terraform import stacklet_repository.example https://github.com/organization/repository
``` 