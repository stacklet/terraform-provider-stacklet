# stacklet_account_discovery (Resource)

Manages an account discovery configuration in Stacklet. Account discovery configurations allow you to automatically discover and manage cloud accounts across different providers.

## Example Usage

```hcl
# Configure AWS account discovery using Organizations
resource "stacklet_account_discovery" "aws" {
  name        = "aws-org-discovery"
  provider    = "aws"
  description = "AWS account discovery via Organizations"
  
  # AWS-specific configuration
  org_read_role   = "arn:aws:iam::123456789012:role/OrganizationReadRole"
  member_role     = "arn:aws:iam::${account}:role/AssetDBRole"
  custodian_role  = "arn:aws:iam::${account}:role/CloudCustodianRole"
  
  suspended = false
}

# Configure Azure account discovery
resource "stacklet_account_discovery" "azure" {
  name        = "azure-discovery"
  provider    = "azure"
  description = "Azure subscription discovery"
  
  # Azure-specific configuration
  tenant_id     = "00000000-0000-0000-0000-000000000000"
  client_id     = "00000000-0000-0000-0000-000000000000"
  client_secret = var.azure_client_secret
  
  suspended = false
}

# Configure GCP account discovery
resource "stacklet_account_discovery" "gcp" {
  name        = "gcp-org-discovery"
  provider    = "gcp"
  description = "GCP project discovery"
  
  # GCP-specific configuration
  org_id          = "1234567890"
  root_folder_ids = ["folder1", "folder2"]
  credential_json = file("path/to/service-account-key.json")
  
  suspended = false
}
```

## Argument Reference

* `name` - (Required) The unique name of the account discovery configuration.
* `provider` - (Required) The cloud provider to discover accounts from (aws, azure, or gcp). This value cannot be changed after creation.
* `description` - (Optional) Human-readable notes about the account discovery configuration.
* `suspended` - (Optional) Whether account discovery is suspended. Defaults to false.

### AWS-specific Arguments

* `org_read_role` - (Required for AWS) The ARN of an IAM role which has permission to read organization data.
* `member_role` - (Optional) Optional IAM role ARN template for AssetDB.
* `custodian_role` - (Optional) Optional IAM role name or ARN template for Cloud Custodian.

### Azure-specific Arguments

* `tenant_id` - (Required for Azure) The Azure tenant ID.
* `client_id` - (Required for Azure) The Azure client ID.
* `client_secret` - (Required for Azure, Sensitive) The Azure client secret.

### GCP-specific Arguments

* `org_id` - (Required for GCP) The GCP organization ID.
* `root_folder_ids` - (Optional) Optional list of GCP folder IDs to scan.
* `exclude_folder_ids` - (Optional) Optional list of GCP folder IDs to exclude from scanning.
* `credential_json` - (Required for GCP, Sensitive) The contents of a JSON-formatted key file for a GCP service account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GraphQL Node ID of the account discovery configuration.
* `validity` - Information about the most recent credential validation attempt:
  ```json
  {
    "valid": true,
    "message": "Successfully validated credentials"
  }
  ```

## Import

Account discovery configurations can be imported using their name:

```shell
terraform import stacklet_account_discovery.aws aws-org-discovery
``` 