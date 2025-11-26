## 0.5.1 - 2025-11-26

- Drop releases for 32bit architectures.
- Update resources and datasources related to MS Teams with new fields.


## 0.5.0 - 2025-11-24

- Add support for dynamic account groups via `dynamic_filter` in `stacklet_account_group` resource/datasource,
- Fix: don't override provider endpoint and api key with the stacklet-admin config if they're declared elsewhere.
- Fix: Don't try logging HTTP responses if the request failed.


## 0.4.0 - 2025-11-19

- Replace `teams_delivery_settings` support in `stacklet_report_group` resource/datasource with `msteams_delivery_settings`, based on the new integration.
- Add support for the following resource types:
  * `stacklet_configuration_profile_email`
  * `stacklet_configuration_profile_jira`
  * `stacklet_configuration_profile_msteams`
  * `stacklet_configuration_profile_servicenow`
  * `stacklet_configuration_profile_slack`
  * `stacklet_configuration_profile_symphony`
- Add support for the following data source types:
  * `stacklet_configuration_profile_msteams`
  * `stacklet_msteams_integration_surface`
- Remove support for the `stacklet_configuration_profile_teams` data source.
- Add missing computed `source` attribute for `stacklet_report_group` resource.


## 0.3.0 - 2025-08-26

- Add module versions for `stacklet_platform` data source.
- Add support for the following resource types:
  * `stacklet_configuration_profile_account_owners`
  * `stacklet_configuration_profile_resource_owner`
  * `stacklet_notification_template`
  * `stacklet_report_group`
- Add support for the following data source types:
  * `stacklet_configuration_profile_account_owners`
  * `stacklet_configuration_profile_email`
  * `stacklet_configuration_profile_jira`
  * `stacklet_configuration_profile_resource_owner`
  * `stacklet_configuration_profile_servicenow`
  * `stacklet_configuration_profile_slack`
  * `stacklet_configuration_profile_symphony`
  * `stacklet_configuration_profile_teams`
  * `stacklet_notification_template`
  * `stacklet_report_group`


## 0.2.0 - 2025-07-02

- Add `stacklet_platform` data source.
- Handle pagination in policy collection API.


## 0.1.0 - 2025-06-17

- First release
- Add support for the following resource types:
  * `stacklet_account`
  * `stacklet_account_discovery_aws`
  * `stacklet_account_discovery_azure`
  * `stacklet_account_discovery_gcp`
  * `stacklet_account_group`
  * `stacklet_account_group_mapping`
  * `stacklet_binding`
  * `stacklet_policy_collection`
  * `stacklet_policy_collection_mapping`
  * `stacklet_repository`
- Add support for the following data source types:
  * `stacklet_account`
  * `stacklet_account_group`
  * `stacklet_binding`
  * `stacklet_policy`
  * `stacklet_policy_collection`
  * `stacklet_repository`
