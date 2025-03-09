# stacklet_account_discovery (Data Source)

Retrieves information about an account discovery configuration in Stacklet. Account discovery configurations allow you to automatically discover and manage cloud accounts across different providers.

## Example Usage

```hcl
# Fetch an account discovery configuration
data "stacklet_account_discovery" "example" {
  name = "aws-org-discovery"
}

# Output discovery details
output "discovery_config" {
  value = data.stacklet_account_discovery.example.config
  sensitive = true
}

output "discovery_validity" {
  value = data.stacklet_account_discovery.example.validity
}
```

## Argument Reference

* `name` - (Required) The unique name of the account discovery configuration.

## Attribute Reference

* `id` - The GraphQL Node ID of the account discovery configuration.
* `description` - Human-readable notes about the account discovery configuration.
* `provider` - The cloud provider to discover accounts from (aws, azure, or gcp).
* `config` - (Sensitive) JSON-encoded configuration specific to the provider:
  * For AWS:
    ```json
    {
      "org_read_role": "arn:aws:iam::123456789012:role/OrganizationReadRole",
      "member_role": "arn:aws:iam::${account}:role/AssetDBRole",
      "custodian_role": "arn:aws:iam::${account}:role/CloudCustodianRole"
    }
    ```
  * For Azure:
    ```json
    {
      "tenant_id": "00000000-0000-0000-0000-000000000000",
      "client_id": "00000000-0000-0000-0000-000000000000",
      "client_secret": "secret"
    }
    ```
  * For GCP:
    ```json
    {
      "org_id": "1234567890",
      "root_folder_ids": ["folder1", "folder2"],
      "exclude_folder_ids": ["folder3"],
      "credential_json": "{...}"
    }
    ```
* `schedule` - JSON-encoded schedule information for when the discovery runs:
  ```json
  {
    "suspended": false
  }
  ```
* `validity` - JSON-encoded information about the most recent credential validation attempt:
  ```json
  {
    "valid": true,
    "message": "Successfully validated credentials"
  }
  ``` 