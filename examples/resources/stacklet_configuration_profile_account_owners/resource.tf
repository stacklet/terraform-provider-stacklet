resource "stacklet_configuration_profile_account_owners" "example" {
  tags = ["Owner", "AccountOwner"]
  default = [
    {
      account = "1111111111"
      owners = [
        "user1@example.com",
        "user2@example.com",
      ]
    },
    {
      account = "2222222222"
      owners = [
        "user3@example.com",
        "user4@example.com",
      ]
    },
  ]
}
