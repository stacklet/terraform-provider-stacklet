# stacklet_policy_collection (Resource)

Manages a policy collection in Stacklet. Policy collections are groups of policies that can be applied to cloud resources, organized by cloud provider.

## Example Usage

```hcl
# Create an AWS policy collection
resource "stacklet_policy_collection" "aws_security" {
  name          = "aws-security-policies"
  cloud_provider = "aws"
  description   = "Security policies for AWS resources"
  auto_update   = true
}

# Create an Azure policy collection
resource "stacklet_policy_collection" "azure_compliance" {
  name          = "azure-compliance-policies"
  cloud_provider = "azure"
  description   = "Compliance policies for Azure resources"
  auto_update   = false
}

# Create a GCP policy collection
resource "stacklet_policy_collection" "gcp_cost" {
  name          = "gcp-cost-policies"
  cloud_provider = "gcp"
  description   = "Cost optimization policies for GCP resources"
  auto_update   = true
}
```

## Argument Reference

* `name` - (Required) The name of the policy collection.
* `cloud_provider` - (Required) The cloud provider for the policy collection (aws, azure, or gcp).
* `description` - (Optional) A description of the policy collection.
* `auto_update` - (Optional) Whether the policy collection automatically updates policy versions. Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GraphQL Node ID of the policy collection.
* `uuid` - The UUID of the policy collection.
* `system` - Whether this is a system policy collection (not user editable).
* `repository` - The repository URL if this collection was created from a repo control file.

## Import

Policy collections can be imported using their UUID:

```shell
terraform import stacklet_policy_collection.aws_security 00000000-0000-0000-0000-000000000000
``` 