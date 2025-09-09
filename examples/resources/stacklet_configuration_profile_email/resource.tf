resource "stacklet_configuration_profile_email" "example_ses" {
  from = "no-reply@example.com"

  ses_region = "us-east-1"
}


resource "stacklet_configuration_profile_email" "example_smtp" {
  from = "no-reply@example.com"

  smtp = {
    server = "smtp.example.com"
    port   = "587"
    ssl    = true

    username            = "myuser@example.com"
    password_wo         = "mysecret"
    password_wo_version = "1"
  }
}
