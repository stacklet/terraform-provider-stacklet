resource "stacklet_account_discovery_aws" "example" {
  name           = "aws"
  description    = "Discovery of AWS account"
  org_read_role  = "arn:aws:iam::123456789012:role/StackletOrgReadRole"
  member_role    = "arn:aws:iam::{Id}:role/StackletAssetDBRole"
  custodian_role = "arn:aws:iam::{id}:role/StackletCustodianRole-{tag}"
}
