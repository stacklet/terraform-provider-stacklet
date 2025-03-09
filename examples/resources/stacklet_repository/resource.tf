resource "stacklet_repository" "example" {
  name        = "example-repo"
  url         = "https://github.com/example/repo"
  description = "Example repository containing Stacklet policies"

  # Optional: Specify which file suffixes to scan for policies
  policy_file_suffix = [".yaml", ".yml", ".json"]

  # Optional: Only scan specific directories for policies
  policy_directories = ["policies", "rules"]

  # Optional: Use a specific branch
  branch_name = "main"

  # Optional: Authentication for private repositories
  auth_user = "git-user"
  auth_token = "ghp_xxxxxxxxxxxxxxxxxxxx"  # or use STACKLET_REPO_TOKEN env var

  # Optional: SSH authentication
  # ssh_private_key = file("~/.ssh/id_rsa")
  # ssh_passphrase = "your-passphrase"

  # Optional: Deep import to scan all commits
  deep_import = true
} 