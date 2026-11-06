resource "stacklet_iac_policy_source" "example" {
  name = "example-public"
  url  = "https://github.com/example/public-repo"
}

resource "stacklet_iac_policy_source" "example_options" {
  name       = "example-options"
  url        = "https://github.com/example/policies-repo"
  branch     = "main"
  paths      = ["policies/terraform", "policies/cloudformation"]
  extensions = [".yaml"]
  auto_scan  = true
}

resource "stacklet_iac_policy_source" "example_http_auth" {
  name = "example-http-auth"
  url  = "https://github.com/example/private-repo"

  token_auth = {
    user             = "some-user"
    token_wo         = "ghp_xxxxxxxxxxxxxxxxxxxx"
    token_wo_version = "1"
  }
}

resource "stacklet_iac_policy_source" "example_ssh_auth" {
  name = "example-ssh-auth"
  url  = "ssh://git@github.com/example/private-repo"

  ssh_auth = {
    user                   = "git"
    private_key_wo         = <<EOT
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1GwAAAJjX7M8q1+zP
KgAAAAtzc2gtZWQyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1Gw
AAAEC/dsi/DISHYy8HxIrX5JWLWhYKv2XFBlL15NLRzIlA5tBIt1S68mWzaTbxsvJqUsKt
ktoEb85tBcGQQmXy8HUbAAAADnJlcG8tdGVzdC1yYXcKAQIDBAUGBw==
-----END OPENSSH PRIVATE KEY-----
EOT
    private_key_wo_version = "1"
  }
}
