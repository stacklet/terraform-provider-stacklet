resource "stacklet_repository" "example_public" {
  name        = "example-public"
  url         = "https://github.com/example/public-repo"
  description = "Example repository containing Stacklet policies"
}

resource "stacklet_repository" "example_http_auth" {
  name = "example-http-auth"
  url  = "https://github.com/example/private-repo"

  auth_user             = "some-user"
  auth_token_wo         = "ghp_xxxxxxxxxxxxxxxxxxxx"
  auth_token_wo_version = 1
}

resource "stacklet_repository" "example_ssh_auth" {
  name = "example-http-auth"
  url  = "ssh://git@github.com/example/private-repo"

  ssh_private_key_wo = <<EOT
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1GwAAAJjX7M8q1+zP
KgAAAAtzc2gtZWQyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1Gw
AAAEC/dsi/DISHYy8HxIrX5JWLWhYKv2XFBlL15NLRzIlA5tBIt1S68mWzaTbxsvJqUsKt
ktoEb85tBcGQQmXy8HUbAAAADnJlcG8tdGVzdC1yYXcKAQIDBAUGBw==
-----END OPENSSH PRIVATE KEY-----
EOT
  ssh_private_key_wo_version = 1
}

resource "stacklet_repository" "example_codecommit" {
  name = "example-codecommit"
  url  = "https://git-codecommit.us-east-1.amazonaws.com/v1/repos/example"

  auth_token_wo         = "arn:aws:iam::123456789012:role/example-repo-access"
  auth_token_wo_version = 1
}
