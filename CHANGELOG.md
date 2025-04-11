## v1.11.0 [2025-04-11]

_What's new?_

- New tables added  
  - [aws_cloudwatch_log_delivery_destination](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_delivery_destination) ([#2469](https://github.com/turbot/steampipe-plugin-aws/pull/2469))  
  - [aws_cloudwatch_log_delivery_source](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_delivery_source) ([#2469](https://github.com/turbot/steampipe-plugin-aws/pull/2469))  
  - [aws_cloudwatch_log_delivery](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_delivery) ([#2469](https://github.com/turbot/steampipe-plugin-aws/pull/2469))  
  - [aws_cloudwatch_log_destination](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_destination) ([#2469](https://github.com/turbot/steampipe-plugin-aws/pull/2469))  
  - [aws_elasticache_update_action](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_update_action) ([#2431](https://github.com/turbot/steampipe-plugin-aws/pull/2431)) (Thanks [@fyqtian](https://github.com/fyqtian) for the contribution!)
  - [aws_quicksight_account_settings](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_account_settings) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_quicksight_data_set](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_data_set) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_quicksight_data_source](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_data_source) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_quicksight_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_group) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_quicksight_namespace](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_namespace) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_quicksight_user](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_user) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_quicksight_vpc_connection](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_quicksight_vpc_connection) ([#2467](https://github.com/turbot/steampipe-plugin-aws/pull/2467))  
  - [aws_s3_multipart_upload](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_s3_multipart_upload) ([#2456](https://github.com/turbot/steampipe-plugin-aws/pull/2456))

_Enhancements_

- Added `folder` metadata to the documentation of all the AWS tables for improved organization on the Steampipe Hub. ([#2465](https://github.com/turbot/steampipe-plugin-aws/pull/2465))  
- Added `inline_policy` and `inline_policy_std` columns to `aws_ssoadmin_permission_set` table. ([#2458](https://github.com/turbot/steampipe-plugin-aws/pull/2458)) (Thanks [@2XXE-SRA](https://github.com/2XXE-SRA) for the contribution!)
- Updated display name to `AWS`.

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/releases/tag/v5.11.5) ([#2460](https://github.com/turbot/steampipe-plugin-aws/pull/2460))  

## v1.10.0 [2025-03-18]

_What's new?_

- New tables added
  - [aws_rds_pending_maintenance_action](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_pending_maintenance_action) ([#2430](https://github.com/turbot/steampipe-plugin-aws/pull/2430)) (Thanks [@fyqtian](https://github.com/fyqtian) for the contribution!)

_Bug fixes_

- Fixed `aws_health_*` tables to correctly reference the AWS Health Global endpoint instead of regional endpoints. ([#2450](https://github.com/turbot/steampipe-plugin-aws/pull/2450))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.11.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5114-2025-03-12).
- Recompiled plugin with `golang.org/x/net` with `v0.36.0`. ([#2447](https://github.com/turbot/steampipe-plugin-aws/pull/2447))

## v1.9.0 [2025-03-07]

_What's new?_

- New tables added
  - [aws_lakeformation_permission](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lakeformation_permission) ([#2417](https://github.com/turbot/steampipe-plugin-aws/pull/2417))
  - [aws_lakeformation_resource](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lakeformation_resource) ([#2417](https://github.com/turbot/steampipe-plugin-aws/pull/2417))
  - [aws_lakeformation_tag](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lakeformation_tag) ([#2417](https://github.com/turbot/steampipe-plugin-aws/pull/2417))

_Enhancements_

- Updated `aws_acm_*`, `aws_sns_*`, `aws_sqs_*`, `aws_cloudtrail_*`, and `aws_guardduty_*` tables to use AWS Go SDK V2, enabling dynamic region listing for all AWS partitions. ([#2440](https://github.com/turbot/steampipe-plugin-aws/pull/2440)) 

## v1.8.0 [2025-02-28]

_Dependencies_

- Recompiled plugin with Go version `1.23.1`.
- Recompiled plugin with [steampipe-plugin-sdk v5.11.3](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5113-2025-02-11) that addresses critical and high vulnerabilities in dependent packages.

## v1.7.0 [2025-02-28]

_Enhancements_

- Added support for `sse_customer_algorithm`, `sse_customer_key` and `sse_customer_key_md5` optional key quals in the `aws_s3_object` table to list objects encrypted with SSE-C. ([#2409](https://github.com/turbot/steampipe-plugin-aws/pull/2409))
- Added parent hydrate support in the `aws_ecr_image_scan_finding` table to manage the complex join queries. ([#2376](https://github.com/turbot/steampipe-plugin-aws/pull/2376))
- Added `pending_modified_values` column to the `aws_rds_db_instance` table. ([#2411](https://github.com/turbot/steampipe-plugin-aws/pull/2411))
- Added tags to `aws_glue_*` tables. ([#2402](https://github.com/turbot/steampipe-plugin-aws/pull/2402)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Added tag retrieval example to `aws_ses_domain_identity` table documentation. ([#2432](https://github.com/turbot/steampipe-plugin-aws/pull/2432))
- Added `logging_config` column to the `aws_lambda_function` table. ([#2423](https://github.com/turbot/steampipe-plugin-aws/pull/2423))

_Bug fixes_

- Fixed the `nil pointer dereference` error when querying AWS RDS custom instances. ([#2436](https://github.com/turbot/steampipe-plugin-aws/pull/2436))
- Fixed the `region` column of `aws_wafregional_rule` table to correctly return the resource region instead of `global`. ([#2429](https://github.com/turbot/steampipe-plugin-aws/pull/2429))
- Fixed the `arn` column in `aws_vpc_eip` table to use the correct format. ([#2415](https://github.com/turbot/steampipe-plugin-aws/pull/2415)) (Thanks [@thomasklemm](https://github.com/thomasklemm) for the contribution!)
- Fixed not found errors in `aws_kinesis_consumer` and `aws_lightsail_instance` tables. ([#2408](https://github.com/turbot/steampipe-plugin-aws/pull/2408))
- Fixed the `InvalidParameterException` error in `aws_ecs_service` tables when listing tags for older ECS services. ([#2410](https://github.com/turbot/steampipe-plugin-aws/pull/2410))

## v1.6.0 [2025-02-06]

_Enhancements_

- Added columns `bootstrap_broker_string` and `bootstrap_broker_string_tls` to the `aws_msk_cluster` table. ([#2390](https://github.com/turbot/steampipe-plugin-aws/pull/2390)) (Thanks [@insummersnow](https://github.com/insummersnow) for the contribution!)
- Added pagination in the `aws_ec2_ami_shared` table. ([#2260](https://github.com/turbot/steampipe-plugin-aws/pull/2260))
- Added columns `owner_ids` and `image_ids` to the `aws_ec2_ami_shared` table. ([#2260](https://github.com/turbot/steampipe-plugin-aws/pull/2260))

## v1.5.0 [2025-01-03]

_What's new?_

- New tables added
  - [aws_costoptimizationhub_recommendation](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_costoptimizationhub_recommendation) ([#2355](https://github.com/turbot/steampipe-plugin-aws/pull/2355))
  - [aws_scheduler_schedule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_scheduler_schedule) ([#2359](https://github.com/turbot/steampipe-plugin-aws/pull/2359))
  - [aws_keyspaces_keyspace](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_keyspaces_keyspace) ([#2271](https://github.com/turbot/steampipe-plugin-aws/pull/2271))
  - [aws_keyspaces_table](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_keyspaces_table) ([#2271](https://github.com/turbot/steampipe-plugin-aws/pull/2271))
  - [aws_config_delivery_channel](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_config_delivery_channel) ([#2343](https://github.com/turbot/steampipe-plugin-aws/pull/2343))

_Enhancements_

- Added `instance_type_pattern` column as an optional qual to the `aws_ec2_instance_type` table. ([#2301](https://github.com/turbot/steampipe-plugin-aws/pull/2301))
- Added `image_digest` column as an optional qual to the `aws_ecr_image_scan_finding` table. ([#2357](https://github.com/turbot/steampipe-plugin-aws/pull/2357))
- Added `created_at` and `updated_at` columns as optional quals to the `aws_securityhub_finding` table. ([#2298](https://github.com/turbot/steampipe-plugin-aws/pull/2298))
- Added `account_password_present` column to `aws_iam_account_summary` table. ([#2346](https://github.com/turbot/steampipe-plugin-aws/pull/2346))
- Add `tags` column to `aws_backup_plan table`. ([#2336](https://github.com/turbot/steampipe-plugin-aws/pull/2336)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Bug fixes_

- Fixed the `aws_rds_db_instance` table to correctly return data instead of an error by ignoring the `CertificateNotFound` error code. ([#2363](https://github.com/turbot/steampipe-plugin-aws/pull/2363))

## v1.4.0 [2024-11-22]

_What's new?_

- New tables added
  - [aws_cost_by_region_monthly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_region_monthly) ([#2310](https://github.com/turbot/steampipe-plugin-aws/pull/2310)) (Thanks [@razbne](https://github.com/razbne) for the contribution!)

_Enhancements_

- Added `error`, `is_public`, `resource_owner_account` and `resource_type` optional quals for `aws_accessanalyzer_finding` table. ([#2331](https://github.com/turbot/steampipe-plugin-aws/pull/2331)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
- Updated the `aws_s3_object` table to use the `HeadObject` API to retrieve object metadata. ([#2312](https://github.com/turbot/steampipe-plugin-aws/pull/2312)) (Thanks [@JonMerlevede](https://github.com/JonMerlevede) for the contribution!)

_Bug fixes_

- Fixed the `aws_s3_bucket` table to correctly return data by ignoring the not found error in `getBucketTagging` and `getBucketWebsite` hydrate functions. ([#2335](https://github.com/turbot/steampipe-plugin-aws/pull/2335))

## v1.3.0 [2024-11-14]

_Enhancements_

- Added `multi_region` and `multi_region_configuration` columns to `aws_kms_key` table. ([#2338](https://github.com/turbot/steampipe-plugin-aws/pull/2338)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Bug fixes_

- Fixed the comparison operator `(<= or >=)` for number and date filter in `aws_inspector2_finding` table. ([#2332](https://github.com/turbot/steampipe-plugin-aws/pull/2332)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)

## v1.2.0 [2024-11-04]

_What's new?_

- New tables added
  - [aws_shield_attack](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_attack) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
  - [aws_shield_attack_statistic](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_attack_statistic) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
  - [aws_shield_drt_access](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_drt_access) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
  - [aws_shield_emergency_contact](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_emergency_contact) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
  - [aws_shield_protection](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_protection) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
  - [aws_shield_protection_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_protection_group) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)
  - [aws_shield_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_shield_subscription) ([#2315](https://github.com/turbot/steampipe-plugin-aws/pull/2315)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)

_Enhancements_

- Added `epss_score` column to `aws_inspector2_finding` table. ([#2321](https://github.com/turbot/steampipe-plugin-aws/pull/2321)) (Thanks [@dbermuehler](https://github.com/dbermuehler) for the contribution!)

_Bug fixes_

- Fixed the `aws_ssm_document_permission` table to correctly return `nil` whenever `InvalidDocument` error is returned by the API. ([#2326](https://github.com/turbot/steampipe-plugin-aws/pull/2326))
- Fixed error handling for `aws_iam_user` and `aws_s3_bucket` tables. ([#2324](https://github.com/turbot/steampipe-plugin-aws/pull/2324)) (Thanks [@danielgrittner](https://github.com/danielgrittner) for the contribution!)
- Updated SQL queries to exclude removed table columns. ([#2328](https://github.com/turbot/steampipe-plugin-aws/pull/2328))

## v1.0.1 [2024-10-25]

_Bug fixes_

- Added `verification_token` column to`aws_ses_domain_identity` table which was accidentally removed in v1.0.0.

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Breaking changes_

- The following tables have had deprecated columns removed: ([#2323](https://github.com/turbot/steampipe-plugin-aws/pull/2323))
  - `aws_accessanalyzer_analyzer`:
    - `findings` (replaced by `aws_accessanalyzer_finding` table)
  - `aws_ecr_repository`:
    - `image_details` (replaced by `aws_ecr_image` table)
    - `image_scanning_findings` (replaced by `aws_ecr_image_scan_finding` table)
  - `aws_ecrpublic_repository`:
    - `image_details` (replaced by `aws_ecr_image` table)
  - `aws_glue_job`:
    - `allocated_capacity` (replaced by `max_capacity` column)
  - `aws_securityhub_finding`:
    - `workflow_state` (replaced by `workflow_status` column)
  - `aws_ses_email_identity`:
    - `verification_token`
  - `aws_ssm_document`:
    - `account_ids` (replaced by `aws_ssm_document_permission` table)
    - `account_sharing_info_list` (replaced by `aws_ssm_document_permission` table)
  - `aws_vpc_security_group_rule`:
    - `cidr_ip` (replaced by `cidr_ipv4` column)
    - `group_name`
    - `owner_id` (replaced by `group_owner_id` column)
    - `pair_group_id` (replaced by `referenced_group_id` column)
    - `pair_group_name`
    - `pair_peering_status` (replaced by `referenced_peering_status` column)
    - `pair_user_id` (replaced by `referenced_user_id` column)
    - `pair_vpc_id` (replaced by `referenced_vpc_id` column)
    - `pair_vpc_peering_connection_id` (replaced by `referenced_vpc_peering_connection_id` column)
    - `vpc_id`

_Enhancements_

- Added `stream_mode_details` column to `aws_kinesis_stream` table. ([#2320](https://github.com/turbot/steampipe-plugin-aws/pull/2320)) (Thanks [@kaushikkishore](https://github.com/kaushikkishore) for the contribution!)

_Bug fixes_

- Fixed the `GetConfig` of the `aws_servicequotas_service_quota` table to correctly return data instead of an error by adding `region` as a required qual. ([#2314](https://github.com/turbot/steampipe-plugin-aws/pull/2314))

## v0.147.0 [2024-09-13]

_Enhancements_

- Added the `event_region` column to the `aws_health_event` table. ([#2293](https://github.com/turbot/steampipe-plugin-aws/pull/2293))
- Added the `location_type` column to the  `aws_ec2_instance_type` table. ([#2294](https://github.com/turbot/steampipe-plugin-aws/pull/2294))

_Bug fixes_

- Removed unnecessary hydration of the `instance_type` column in `aws_ec2_instance_type` table. ([#2294](https://github.com/turbot/steampipe-plugin-aws/pull/2294))
- Fixed an issue where credentials from import foreign schema were lost after restarting session in the Posgres FDW extensions of the plugin. ([#2275](https://github.com/turbot/steampipe-plugin-aws/issues/2275))

## v0.146.0 [2024-09-03]

_What's new?_

- New tables added
  - [aws_app_runner_service](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_app_runner_service) ([#2279](https://github.com/turbot/steampipe-plugin-aws/pull/2279))
  - [aws_ec2_load_balancer_listener_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_load_balancer_listener_rule) ([#2272](https://github.com/turbot/steampipe-plugin-aws/pull/2272))
  - [aws_lightsail_bucket](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lightsail_bucket) ([#2258](https://github.com/turbot/steampipe-plugin-aws/pull/2258))
  - [aws_memorydb_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_memorydb_cluster) ([#2268](https://github.com/turbot/steampipe-plugin-aws/pull/2268))
  - [aws_securityhub_enabled_product_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_enabled_product_subscription) ([#2281](https://github.com/turbot/steampipe-plugin-aws/pull/2281))
  - [aws_timestreamwrite_database](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_timestreamwrite_database) ([#2269](https://github.com/turbot/steampipe-plugin-aws/pull/2269))
  - [aws_timestreamwrite_table](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_timestreamwrite_table) ([#2269](https://github.com/turbot/steampipe-plugin-aws/pull/2269))

_Enhancements_

- Updated the `aws_ec2_ami` table to correctly return disabled AMIs on passing `disabled` value to the `state` optional qual (`where state = 'disabled'`). ([#2277](https://github.com/turbot/steampipe-plugin-aws/pull/2277))

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#2283](https://github.com/turbot/steampipe-plugin-aws/pull/2283))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#2286](https://github.com/turbot/steampipe-plugin-aws/pull/2286))

## v0.145.0 [2024-08-07]

_Enhancements_

- Added `location_type` column as an optional qual to the `aws_ec2_instance_availability` table and 6 new columns to the `aws_ec2_instance_type` table. ([#2078](https://github.com/turbot/steampipe-plugin-aws/pull/2078))
- Updated docs for `aws_appautoscaling_policy` and `aws_appautoscaling_target` tables to add information on required quals. ([#2247](https://github.com/turbot/steampipe-plugin-aws/pull/2247))
- Added the `type` column as an optional qual to the `aws_auditmanager_control` table. ([#2254](https://github.com/turbot/steampipe-plugin-aws/pull/2254))

_Bug fixes_

- Fixed the `GetConfig` definition of the `aws_auditmanager_control` table to correctly return data instead of an error. ([#2254](https://github.com/turbot/steampipe-plugin-aws/pull/2254))
- Fixed the `aws_kms_key_rotation` table to correctly return `nil` whenever an `AccessDeniedException` error is returned by the API. ([#2253](https://github.com/turbot/steampipe-plugin-aws/pull/2253))

## v0.144.0 [2024-07-10]

_Enhancements_

- Updated IAM parliament permissions to the latest. ([#2239](https://github.com/turbot/steampipe-plugin-aws/pull/2239))

_Bug fixes_

- Fixed the caching issue in 29 tables to correctly return data by adding the missing `CacheMatch: query_cache.CacheMatchExact` property. ([#2232](https://github.com/turbot/steampipe-plugin-aws/pull/2232))
- Fixed the `user_data` column of `aws_ec2_instance` table to remove invalid UTF-8 characters. ([#2240](https://github.com/turbot/steampipe-plugin-aws/pull/2240))

## v0.143.0 [2024-07-05]

_What's new?_

- New tables added
  - [aws_rds_db_recommendation](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_recommendation) ([#2238](https://github.com/turbot/steampipe-plugin-aws/pull/2238))

_Bug fixes_

- Fixed the caching issue in `aws_organizations_account` table. ([#2236](https://github.com/turbot/steampipe-plugin-aws/pull/2236))
- Fixed typo (missing comma) in an example query of `aws_health_affected_entity` table doc. ([#2237](https://github.com/turbot/steampipe-plugin-aws/pull/2237)) (Thanks [@tieum](https://github.com/tieum) for the contribution!)

## v0.142.0 [2024-07-04]

_Enhancements_

- Added 16 new columns to the `aws_lambda_version` table. ([#2229](https://github.com/turbot/steampipe-plugin-aws/pull/2229))

_Bug fixes_

- Fixed the export tool of the plugin to return a non-zero error code instead of 0 whenever an error occurred. ([#79](https://github.com/turbot/steampipe-export/pull/79))

## v0.141.0 [2024-07-01]

_Bug fixes_

- Reverted the Export CLI behaviour to return `<nil>` for `null` values instead of `empty`. ([#77](https://github.com/turbot/steampipe-export/issues/77))

## v0.140.0 [2024-06-28]

_What's new_

- New tables added
  - [aws_codestar_notification_rule](//hub.steampipe.io/plugins/turbot/aws/tables/aws_codestar_notification_rule) ([#2217](https://github.com/turbot/steampipe-plugin-aws/pull/2217))

_Enhancements_

- Added 9 new columns to the `aws_elasticache_cluster` table. ([#2224](https://github.com/turbot/steampipe-plugin-aws/pull/2224))

_Bug fixes_

- Fixed the `aws_s3_object` table not returning any rows due to panic error. ([#2221](https://github.com/turbot/steampipe-plugin-aws/pull/2221))
- Fixed no rows being returned from the `aws_organizations_account` table if an unqualified query is run before one with `parent_id` specified.
- Fixed data type for `configuration_endpoint` column in `aws_elasticache_cluster` table to be `json`. ([#2214](https://github.com/turbot/steampipe-plugin-aws/pull/2214))

## v0.139.0 [2024-06-17]

_What's new?_

- New tables added
  - [aws_route53_vpc_association_authorization](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_vpc_association_authorization) ([#2199](https://github.com/turbot/steampipe-plugin-aws/pull/2199)) (Thanks [@jramosf](https://github.com/jramosf) for the contribution!)

_Enhancements_

- Updated `aws_s3_bucket`, `aws_s3_bucket_intelligent_tiering_configuration`, `aws_s3_object` and `aws_s3_object_version` tables to use `HeadBucket` API instead of `GetBucketLocation` to fetch the region that the bucket resides in. ([#2082](https://github.com/turbot/steampipe-plugin-aws/pull/2082)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Added column `create_time` to `aws_ec2_key_pair` table. ([#2196](https://github.com/turbot/steampipe-plugin-aws/pull/2196)) (Thanks [@kasadaamos](https://github.com/https://github.com/kasadaamos) for the contribution!)
- Added `instance_type` column as an optional qual to the `aws_ec2_instance_type` table. ([#2200](https://github.com/turbot/steampipe-plugin-aws/pull/2200))

_Bug fixes_

- Fixed the `akas` column in `aws_health_affected_entity` table to correctly return data instead of an error by handling events that do not have any `ARN`. ([#2189](https://github.com/turbot/steampipe-plugin-aws/pull/2189))
- Fixed `cname` and `endpoint_url` columns of `aws_elastic_beanstalk_environment` table to correctly return data instead of `null`. ([#2201](https://github.com/turbot/steampipe-plugin-aws/pull/2201))
- Fixed the `aws_api_gatewayv2_*` tables to correctly return data instead of an error by excluding support for the new `il-central-1` region. ([#2190](https://github.com/turbot/steampipe-plugin-aws/pull/2190))

## v0.138.0 [2024-05-09]

_Enhancements_

- The Plugin and the Steampipe Anywhere binaries are now built with the `netgo` package for both the Linux and Darwin systems. ([#219](https://github.com/turbot/steampipe-plugin-kubernetes/pull/219)) ([#2180](https://github.com/turbot/steampipe-plugin-aws/pull/2180))

_Bug fixes_

- Fixed the `aws_ebs_snapshot` table to correctly return data instead of an empty row. ([#2185](https://github.com/turbot/steampipe-plugin-aws/pull/2185))

_Dependencies_

- Recompiled plugin with [github.com/hashicorp/go-getter v1.7.4](https://github.com/hashicorp/go-getter/releases/tag/v1.7.4). ([#2178](https://github.com/turbot/steampipe-plugin-aws/pull/2178))

## v0.137.0 [2024-04-29]

_What's new?_

- New tables added
  - [aws_iot_thing_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iot_thing_group) ([#1998](https://github.com/turbot/steampipe-plugin-aws/pull/1998))
  - [aws_iot_thing_type](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iot_thing_type) ([#1999](https://github.com/turbot/steampipe-plugin-aws/pull/1999))
  - [aws_kms_key_rotation](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_kms_key_rotation) ([#2174](https://github.com/turbot/steampipe-plugin-aws/pull/2174))

_Enhancements_

- Added the `version` flag to the plugin's Export tool. ([#65](https://github.com/turbot/steampipe-export/pull/65))

_Bug fixes_

- Fixed the broken Postgres 14, Postgres 15 and SQLite x86_64 binaries for Darwin operating systems.
- Fixed intermittent FDW crashes when certain postgres errors resulted in a signal 16 being raised. ([#455](https://github.com/turbot/steampipe-postgres-fdw/pull/455))

## v0.136.1 [2024-04-23]

_Bug fixes_

- Fixed the [runtime error](https://github.com/turbot/steampipe-postgres-fdw/issues/454) in the `v0.136.0` version of the plugin’s Postgres FDW extension.

## v0.136.0 [2024-04-19]

_What's new?_

- New tables added
  - [aws_iot_fleet_metric](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iot_fleet_metric) ([#2000](https://github.com/turbot/steampipe-plugin-aws/pull/2000))

_Enhancements_

- The `account_id` column has now been assigned as a connection key column across all the tables which facilitates more precise and efficient querying across multiple AWS accounts. ([#2133](https://github.com/turbot/steampipe-plugin-aws/pull/2133))

_Bug fixes_

- Fixed the `getDirectoryServiceSnapshotLimit` and `getDirectoryServiceEventTopics` hydrate calls in the `aws_directory_service_directory` table to correctly return `nil` for the unsupported `ADConnector` services instead of an error. ([#2170](https://github.com/turbot/steampipe-plugin-aws/pull/2170))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.10.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v5100-2024-04-10) that adds support for connection key columns. ([#2133](https://github.com/turbot/steampipe-plugin-aws/pull/2133))
- Recompiled plugin with [aws-sdk-go v1.26.1](https://github.com/aws/aws-sdk-go-v2/blob/main/CHANGELOG.md). ([#2163](https://github.com/turbot/steampipe-plugin-aws/pull/2163))

## v0.135.0 [2024-04-12]

_What's new?_

- New tables added
  - [aws_accessanalyzer_finding](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_accessanalyzer_finding) ([#2142](https://github.com/turbot/steampipe-plugin-aws/pull/2142))

_Enhancements_

- Added `snapshot_block_public_access_state` column to `aws_ec2_regional_settings` table. ([#2077](https://github.com/turbot/steampipe-plugin-aws/pull/2077))

_Bug fixes_

- Fixed the `getDirectoryServiceSnapshotLimit` and `getDirectoryServiceEventTopics` hydrate calls in the `aws_directory_service_directory` table to correctly return `nil` for unsupported `SharedMicrosoftAD` services instead of an error. ([#2156](https://github.com/turbot/steampipe-plugin-aws/pull/2156))
- Fixed the plugin's Postgres FDW Extension crash [issue](https://github.com/turbot/steampipe-postgres-fdw/issues/434).

## v0.134.0 [2024-03-29]

_What's new?_

- New tables added
  - [aws_backup_job](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_job) ([#2145](https://github.com/turbot/steampipe-plugin-aws/pull/2145)) (Thanks [@rogerioacp](https://github.com/rogerioacp) for the contribution!)
  - [aws_elastic_beanstalk_application_version](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elastic_beanstalk_application_version) ([#2150](https://github.com/turbot/steampipe-plugin-aws/pull/2150))
  - [aws_rds_db_engine_version](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_engine_version) ([#2098](https://github.com/turbot/steampipe-plugin-aws/pull/2098))
  - [aws_s3_object_version](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_s3_object_version) ([#2070](https://github.com/turbot/steampipe-plugin-aws/pull/2070))
  - [aws_servicequotas_service](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicequotas_service) ([#2070](https://github.com/turbot/steampipe-plugin-aws/pull/2141))

_Enhancements_

- The plugin level logs have been updated to maintain consistency: `Trace` logs have been elevated to `Debug`, `Info` logs elevated to `Error` where needed, and unnecessary `Debug` logs removed to streamline and optimize logging. ([#2131](https://github.com/turbot/steampipe-plugin-aws/pull/2131))

_Bug fixes_

- Fixed the `aws_vpc_eip` table to return an `Access Denied` error instead of an `Invalid Memory Address or Nil Pointer Dereference` error when a `Service Control Policy` is applied to an account for a specific region. ([#2136](https://github.com/turbot/steampipe-plugin-aws/pull/2136))
- Fixed the `aws_s3_bucket` terraform script to prevent the `AccessControlListNotSupported: The bucket does not allow ACLs` error during the `PutBucketAcl` terraform call. ([#2080](https://github.com/turbot/steampipe-plugin-aws/pull/2080)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Fixed an issue where querying regional tables while using AWS profiles with `cross-account` role credentials results in the correct error being reported instead of zero rows. ([#2137](https://github.com/turbot/steampipe-plugin-aws/pull/2137))
- Fixed pagination in the `aws_ebs_snapshot` table to make fewer API calls when the `limit` parameter is passed to the query. ([#2088](https://github.com/turbot/steampipe-plugin-aws/pull/2088))

## v0.133.0 [2024-03-15]

_What's new?_

- New tables added
  - [aws_acmpca_certificate_authority](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_acmpca_certificate_authority) ([#2125](https://github.com/turbot/steampipe-plugin-aws/pull/2125))
  - [aws_dms_endpoint](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dms_endpoint) ([#1992](https://github.com/turbot/steampipe-plugin-aws/pull/1992))
  - [aws_dms_replication_task](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dms_replication_task) ([#2110](https://github.com/turbot/steampipe-plugin-aws/pull/2110))
  - [aws_docdb_cluster_snapshot](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_docdb_cluster_snapshot) ([#2123](https://github.com/turbot/steampipe-plugin-aws/pull/2123))
  - [aws_transfer_user](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_transfer_user) ([#2089](https://github.com/turbot/steampipe-plugin-aws/pull/2089)) (Thanks [@jramosf](https://github.com/jramosf) for the contribution!)

_Enhancements_

- Added `auto_minor_version_upgrade` column to `aws_rds_db_cluster` table. ([#2109](https://github.com/turbot/steampipe-plugin-aws/pull/2109))
- Added `open_zfs_configuration` column to `aws_fsx_file_system` table. ([#2113](https://github.com/turbot/steampipe-plugin-aws/pull/2113))
- Added `logging_configuration` column to `aws_networkfirewall_firewall` table. ([#2115](https://github.com/turbot/steampipe-plugin-aws/pull/2115))
- Added `lf_tags` column to `aws_glue_catalog_table` table. ([#2128](https://github.com/turbot/steampipe-plugin-aws/pull/2128))

_Bug fixes_

- Fixed the query in the `aws_s3_bucket` table doc to correctly filter out buckets without the `application` tag. ([#2093](https://github.com/turbot/steampipe-plugin-aws/pull/2093))
- Fixed the `aws_cloudtrail_lookup_event` input param to pass correctly `end_time` as an optional qual. ([#2102](https://github.com/turbot/steampipe-plugin-aws/pull/2102))
- Fixed the `arn` column of the `aws_elastic_beanstalk_environment` table to correctly return data instead of `null`. ([#2105](https://github.com/turbot/steampipe-plugin-aws/issues/2105))
- Fixed the `template_body_json` column of the `aws_cloudformation_stack` table to correctly return data by adding a new transform function `formatJsonBody`, replacing the `UnmarshalYAML` transform function. ([#1959](https://github.com/turbot/steampipe-plugin-aws/pull/1959))
- Fixed the `next_execution_time` column of `aws_ssm_maintenance_window` table to be of `String` datatype instead of `TIMESTAMP`. ([#2116](https://github.com/turbot/steampipe-plugin-aws/pull/2116))
- Renamed the `client_log_options` column to `connection_log_options` in  `aws_ec2_client_vpn_endpoint` table to correctly return data instead of `null`. ([#2122](https://github.com/turbot/steampipe-plugin-aws/pull/2122))

## v0.132.0 [2024-02-27]

_What's new?_

- New tables added
  - [aws_ecr_registry_scanning_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecr_registry_scanning_configuration) ([#2084](https://github.com/turbot/steampipe-plugin-aws/pull/2084))

_Bug fixes_

- Fixed the `InvalidParameterCombination` error when querying the `aws_rds_db_instance` table. ([#2085](https://github.com/turbot/steampipe-plugin-aws/pull/2085))
- Fixed `aws_rds_db_instance_metric_write_iops_daily` table to correctly display `WriteIOPS` instead of `ReadIOPS`. ([#2079](https://github.com/turbot/steampipe-plugin-aws/pull/2079))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.9.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v590-2024-02-26) that fixes critical caching issues. ([#2067](https://github.com/turbot/steampipe-plugin-aws/pull/2067))

## v0.131.0 [2024-02-15]

_What's new?_

- New tables added
  - [aws_api_gateway_method](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_api_gateway_method) ([#1995](https://github.com/turbot/steampipe-plugin-aws/pull/1995))
  - [aws_dms_certificate](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_organizational_unit) ([#1985](https://github.com/turbot/steampipe-plugin-aws/pull/1985))
  - [aws_emr_security_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_security_configuration) ([#1984](https://github.com/turbot/steampipe-plugin-aws/pull/1984))
  - [aws_iot_thing](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iot_thing) ([#1997](https://github.com/turbot/steampipe-plugin-aws/pull/1997))
  - [aws_organizations_organizational_unit](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_organizational_unit) ([#2063](https://github.com/turbot/steampipe-plugin-aws/pull/2063))
  - [aws_ssmincidents_response_plan](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssmincidents_response_plan) ([#1885](https://github.com/turbot/steampipe-plugin-aws/pull/1885))
  - [aws_trusted_advisor_check_summary](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_trusted_advisor_check_summary) ([#1978](https://github.com/turbot/steampipe-plugin-aws/pull/1978))

_Bug fixes_

- Fixed `aws_sfn_state_machine_execution_history` table to handle pagination and ignore errors for expired executions history. ([#1934](https://github.com/turbot/steampipe-plugin-aws/pull/1934)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Fixed the `aws_health_affected_entity` table to correctly return data instead of an interface conversion error. ([#2072](https://github.com/turbot/steampipe-plugin-aws/pull/2072))

## v0.130.0 [2024-02-02]

_Enhancements_

- Optimized `aws_cloudwatch_log_stream` table's query performance by adding `descending`, `log_group_name`, `log_stream_name_prefix` and `order_by` new optional key qual columns. ([#1951](https://github.com/turbot/steampipe-plugin-aws/pull/1951))
- Optimized `aws_ssm_inventory` table's query performance by adding new optional key qual columns such as `filter_key`, `filter_value`, `network_attribute_key`, `network_attribute_value`, etc. ([#1980](https://github.com/turbot/steampipe-plugin-aws/pull/1980))

_Bug fixes_

- Fixed `aws_cloudwatch_log_group` table key column to be globally unique by filtering the results by region. ([#1976](https://github.com/turbot/steampipe-plugin-aws/pull/1976))
- Removed duplicate memoizing of getCommonColumns function from `aws_s3_multi_region_access_point` and `aws_ec2_launch_template` tables.([#2065](https://github.com/turbot/steampipe-plugin-aws/pull/2065))
- Fixed error for column `type_name` in table `aws_ssm_inventory_entry`. ([#1980](https://github.com/turbot/steampipe-plugin-aws/pull/1980))
- Added the missing rate-limiter tags for `aws_s3_bucket` table's `GetBucketLocation` hydrate function to optimize query performance. ([#2066](https://github.com/turbot/steampipe-plugin-aws/pull/2066))

## v0.129.0 [2024-01-19]

_What's new?_

- New tables added
  - [aws_servicecatalog_provisioned_product](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicecatalog_provisioned_product) ([#1917](https://github.com/turbot/steampipe-plugin-aws/pull/1917))

_Enhancements_

- Added `deletion_protection_enabled` column to `aws_dynamodb_table` table. ([#2049](https://github.com/turbot/steampipe-plugin-aws/pull/2049))

_Bug fixes_

- Fixed default page size in `aws_organizations_account` table. ([#2058](https://github.com/turbot/steampipe-plugin-aws/pull/2058))
- Fixed `processor_features` column in `aws_rds_db_instance` not returning data when default value is set. ([#2028](https://github.com/turbot/steampipe-plugin-aws/pull/2028))
- Temporarily removed `aws_organizations_organizational_unit` table due to LTREE column issue. ([#2058](https://github.com/turbot/steampipe-plugin-aws/pull/2058))

## v0.128.0 [2024-01-15]

_What's new?_

- New tables added
  - [aws_cloudtrail_lookup_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_lookup_event) ([#2047](https://github.com/turbot/steampipe-plugin-aws/pull/2047))
  - [aws_organizations_organizational_unit](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_organizational_unit) ([#1677](https://github.com/turbot/steampipe-plugin-aws/pull/1677))
  - [aws_organizations_root](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_root) ([#1677](https://github.com/turbot/steampipe-plugin-aws/pull/1677))
  - [aws_sns_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sns_subscription) ([#2046](https://github.com/turbot/steampipe-plugin-aws/pull/2046))

**Note :** Table `aws_sns_topic_subscription` will be changing behaviours in a future release to return results from `ListSubscriptionsByTopic` instead of `ListSubscriptions`.

## v0.127.0 [2024-01-10]

_What's new?_

- New tables added
  - [aws_appsync_graphql_api](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_appsync_graphql_api) ([#2027](https://github.com/turbot/steampipe-plugin-aws/pull/2027))
  - [aws_mq_broker](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_mq_broker) ([#2020](https://github.com/turbot/steampipe-plugin-aws/pull/2020))

_Enhancements_

- Added `storage_throughput` column to `aws_rds_db_instance` table. ([#2010](https://github.com/turbot/steampipe-plugin-aws/pull/2010)) (Thanks [@toddwh50](https://github.com/toddwh50) for the contribution!)
- Added `layers` column to `aws_lambda_function` table. ([#2008](https://github.com/turbot/steampipe-plugin-aws/pull/2008)) (Thanks [@icaliskanoglu](https://github.com/icaliskanoglu) for the contribution!)
- Added `tags` column to `aws_backup_recovery_point` and `aws_backup_vault` tables.  ([#2033](https://github.com/turbot/steampipe-plugin-aws/pull/2033))

_Bug fixes_

- Custom HTTP client should allow buildable settings through env var options such as AWS_CA_BUNDLE. ([#2044](https://github.com/turbot/steampipe-plugin-aws/pull/2044))
- Fixed `MaxItems` in `aws_iam_policy` and `aws_iam_policy_attachment` tables to use `1000` instead of `100` to avoid unnecessary API calls. ([#2025](https://github.com/turbot/steampipe-plugin-aws/pull/2025)) ([#2026](https://github.com/turbot/steampipe-plugin-aws/pull/2026))

## v0.126.0 [2023-12-29]

_Enhancements_

- Updated the plugin to use a shared, optimized HTTP client that enhances DNS management and reduces connection floods for more stable and efficient queries. ([#2036](https://github.com/turbot/steampipe-plugin-aws/pull/2036))

## v0.125.0 [2023-12-20]

_Enhancements_

- Updated the `.goreleaser` file to build the netgo package only for Darwin systems. ([#2029](https://github.com/turbot/steampipe-plugin-aws/pull/2029))

## v0.124.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview).
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension.
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/LICENSE).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server enacapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#2011](https://github.com/turbot/steampipe-plugin-aws/pull/2011))

## v0.123.0 [2023-11-16]

_What's new?_

- New tables added
  - [aws_lambda_event_source_mapping](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lambda_event_source_mapping) ([#1874](https://github.com/turbot/steampipe-plugin-aws/issues/1874)) (Thanks [@nickman](https://github.com/nickman) for the contribution!)

_Enhancements_

- Added the `resource_record_set_limit` column to `aws_route53_zone` table. ([#1969](https://github.com/turbot/steampipe-plugin-aws/pull/1969)) (Thanks [@keyolk](https://github.com/keyolk) for the contribution!)

## v0.122.0 [2023-11-10]

_What's new?_

- New tables added
  - [aws_fms_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_fms_policy) ([#1851](https://github.com/turbot/steampipe-plugin-aws/pull/1851))
  - [aws_fms_app_list](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_fms_app_list) ([#1851](https://github.com/turbot/steampipe-plugin-aws/pull/1851))
  - [aws_transfer_server](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_transfer_server) ([#1909](https://github.com/turbot/steampipe-plugin-aws/pull/1909)) (Thanks [@jramosf](https://github.com/jramosf) for the contribution!)

_Enhancements_

- Added the `features` column to `aws_guardduty_detector` table. ([#1958](https://github.com/turbot/steampipe-plugin-aws/pull/1958))

## v0.121.1 [2023-11-06]

_Bug fixes_

- Fixed the description of the `name` column in `aws_organizations_account` table. ([#1947](https://github.com/turbot/steampipe-plugin-aws/pull/1947)) (Thanks [@badideasforsale](https://github.com/badideasforsale) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v563-2023-11-06) which addresses the issue of expired credentials being intermittently retained in the connection cache. ([#1956](https://github.com/turbot/steampipe-plugin-aws/pull/1956))

## v0.121.0 [2023-10-13]

_Enhancements_

- Improved documentation and descriptions for the `aws_iam_role` table. ([#1940](https://github.com/turbot/steampipe-plugin-aws/pull/1940))
- Replaced uses of `rand.Seed` with latest `rand.NewSource`. ([#1933](https://github.com/turbot/steampipe-plugin-aws/pull/1933))

## v0.120.2 [2023-10-04]

_Bug fixes_

- Removed custom plugin level retryer which was unnecessary as the plugin already uses the AWS SDK retryer. ([#1932](https://github.com/turbot/steampipe-plugin-aws/pull/1932))
- The plugin now retries errors with the error code `UnknownError`. These are often thrown by services like SNS when performing a large number of requests. ([#1932](https://github.com/turbot/steampipe-plugin-aws/pull/1932))

## v0.120.1 [2023-10-03]

_Bug fixes_

- Fixed the `source_account_id` column of `aws_securityhub_finding` table to correctly return data instead of `null`. ([#1927](https://github.com/turbot/steampipe-plugin-aws/pull/1927)) (Thanks [@gabrielsoltz](https://github.com/gabrielsoltz) for the contribution!)
- Fixed the `members` column of `aws_rds_db_cluster` table to correctly return data instead of `null`. ([#1926](https://github.com/turbot/steampipe-plugin-aws/pull/1926))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#1930](https://github.com/turbot/steampipe-plugin-aws/pull/1930))

## v0.120.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#1905](https://github.com/turbot/steampipe-plugin-aws/pull/1905))
- Recompiled plugin with Go version `1.21`. ([#1905](https://github.com/turbot/steampipe-plugin-aws/pull/1905))

## v0.119.0 [2023-09-29]

_Enhancements_

- Updated the `Makefile` to build the netgo package only for Darwin systems. ([#1918](https://github.com/turbot/steampipe-plugin-aws/pull/1918))
- Added the `configuration_settings` column to `aws_elastic_beanstalk_environment` table. ([#1916](https://github.com/turbot/steampipe-plugin-aws/pull/1916))

_Bug fixes_

- Fixed the table `aws_dynamodb_backup` to return nil instead of an error when backup does not exist. ([#1914](https://github.com/turbot/steampipe-plugin-aws/pull/1914))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v552-2023-09-29) which improves logging for connection config updates. ([#1921](https://github.com/turbot/steampipe-plugin-aws/pull/1921))

## v0.118.1 [2023-09-14]

_Bug fixes_

- Fixed the data type of `capacity_reservation_specification` column of `aws_ec2_instance` table to be of `JSON` type instead of `STRING`. ([#1903](https://github.com/turbot/steampipe-plugin-aws/pull/1903))

## v0.118.0 [2023-09-07]

_What's new?_

- New tables added
  - [aws_workspaces_directory](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_workspaces_directory) ([#1884](https://github.com/turbot/steampipe-plugin-aws/pull/1884))

_Enhancements_

- Added an example query in the `aws_ec2_instance` table doc for fetching subnet details of instances. ([#1883](https://github.com/turbot/steampipe-plugin-aws/pull/1883)) (Thanks [@Pankaj-SinghR](https://github.com/Pankaj-SinghR) for the contribution!)

_Bug fixes_

- Fixed the data type of the `sms_configuration_failure` column in the `aws_cognito_user_pool` table to be of `STRING` type instead of `JSON`. ([#1890](https://github.com/turbot/steampipe-plugin-aws/pull/1890)) (Thanks [@KTamas](https://github.com/KTamas) for the contribution!)
- Fixed typo in the `listQueryRegionsForConnection` function in the `multi_region.go` file. ([#1887](https://github.com/turbot/steampipe-plugin-aws/pull/1887)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Dependencies_

- Recompiled plugin with `golang.org/x/net v0.7.0`. ([#1864](https://github.com/turbot/steampipe-plugin-aws/pull/1864))

## v0.117.0 [2023-08-25]

_What's new?_

- New tables added
  - [aws_cognito_identity_pool](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cognito_identity_pool) ([#1876](https://github.com/turbot/steampipe-plugin-aws/pull/1876)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Enhancements_

- Added the `engine_type` and `endpoints` columns to `aws_elasticsearch_domain` table. ([#1858](https://github.com/turbot/steampipe-plugin-aws/pull/1858)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v0.116.0 [2023-08-17]

_What's new?_

- New tables added
  - [aws_cognito_identity_provider](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cognito_identity_provider) ([#1854](https://github.com/turbot/steampipe-plugin-aws/pull/1854)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
  - [aws_cognito_user_pool](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cognito_user_pool) ([#1854](https://github.com/turbot/steampipe-plugin-aws/pull/1854)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v0.115.0 [2023-08-08]

_Enhancements_

- Updated the `Makefile` to build plugin in `STEAMPIPE_INSTALL_DIR` if set. ([#1857](https://github.com/turbot/steampipe-plugin-aws/pull/1857)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Added column `offering_class` to `aws_pricing_product` table ([#1863](https://github.com/turbot/steampipe-plugin-aws/pull/1863)) (Thanks [@rasta-rocket](https://github.com/rasta-rocket) for the contribution!)

_Bug fixes_

- Fixed the `aws_ec2_network_load_balancer` table doc to remove the incorrect security group association example. ([#1869](https://github.com/turbot/steampipe-plugin-aws/pull/1869)) (Thanks [@
tinder-tder](https://github.com/tinder-tder) for the contribution!)
- Fixed `aws_rds_db_cluster`, `aws_rds_db_cluster_snapshot`, `aws_rds_db_instance`, `aws_rds_db_snapshot` tables to correctly filter out the `DocDB` and `Neptune` resources. ([#1868](https://github.com/turbot/steampipe-plugin-aws/pull/1868))

## v0.114.0 [2023-08-04]

_What's new?_

- New tables added
  - [aws_neptune_db_cluster_snapshot](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_neptune_db_cluster_snapshot) ([#1866](https://github.com/turbot/steampipe-plugin-aws/pull/1866))

## v0.113.0 [2023-07-28]

_What's new?_

- New tables added
  - [aws_directory_service_log_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_directory_service_log_subscription) ([#1852](https://github.com/turbot/steampipe-plugin-aws/pull/1852))

_Enhancements_

- Added the `fifo_throughput_limit` and `deduplication_scope` columns to the `aws_sqs_queue` table. ([#1859](https://github.com/turbot/steampipe-plugin-aws/pull/1859)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Added the `description` column to the `aws_api_gatewayv2_api` table. ([#1856](https://github.com/turbot/steampipe-plugin-aws/pull/1856)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v0.112.0 [2023-07-20]

_Breaking changes_

- Fixed the `aws_rds_db_*` tables to list out `AWS RDS` resources excluding the `AWS DocDB` ones. Please use `aws_docdb_*` tables instead. ([#1768](https://github.com/turbot/steampipe-plugin-aws/pull/1768))

_What's new?_

- New tables added
  - [aws_directory_service_certificate](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_directory_service_certificate) ([#1836](https://github.com/turbot/steampipe-plugin-aws/pull/1836))
  - [aws_vpc_nat_gateway_metric_bytes_out_to_destination](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_nat_gateway_metric_bytes_out_to_destination) ([#1842](https://github.com/turbot/steampipe-plugin-aws/pull/1842))

_Bug fixes_

- Fixed the optional quals of the `aws_inspector2_finding` table to correctly return data instead of an empty row. ([#1847](https://github.com/turbot/steampipe-plugin-aws/pull/1847))
- Fixed typo in the `aws_vpc_nat_gateway` table doc. ([#1848](https://github.com/turbot/steampipe-plugin-aws/pull/1848)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v0.111.0 [2023-07-14]

_What's new?_

- New tables added
  - [aws_backup_report_plan](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_report_plan) ([#1839](https://github.com/turbot/steampipe-plugin-aws/pull/1839))
  - [aws_cloudformation_stack_set](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudformation_stack_set) ([#1826](https://github.com/turbot/steampipe-plugin-aws/pull/1826))

## v0.110.0 [2023-07-13]

_What's new?_

- New tables added
  - [aws_appstream_fleet](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_appstream_fleet) ([#1838](https://github.com/turbot/steampipe-plugin-aws/pull/1838))

_Enhancements_

- Added the `event_topics` and `snapshot_limit` columns to the `aws_directory_service_directory` table. ([#1833](https://github.com/turbot/steampipe-plugin-aws/pull/1833))

_Bug fixes_

- Fixed the `aws_dlm_lifecycle_policy` table to correctly return results instead of an error. ([#1834](https://github.com/turbot/steampipe-plugin-aws/pull/1834))

## v0.109.1 [2023-07-10]

_Bug fixes_

- Fixed the `certificate` and `certificate_chain` columns of the `aws_acm_certificate` table to correctly return data instead of returning an error. ([#1827](https://github.com/turbot/steampipe-plugin-aws/pull/1827))

## v0.109.0 [2023-07-06]

_What's new?_

- New tables added
  - [aws_iam_open_id_connect_provider](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iam_open_id_connect_provider) ([#1798](https://github.com/turbot/steampipe-plugin-aws/pull/1798)) (Thanks [@LalitLab](https://github.com/LalitLab) for the contribution!)

_Bug fixes_

- Fixed the `aws_route53_record` table to remove the need of passing `zone_id` in the `where` clause, to avoid cross-account access denied errors. ([#1799](https://github.com/turbot/steampipe-plugin-aws/pull/1799))

## v0.108.0 [2023-06-30]

_What's new?_

- New tables added
  - [aws_appautoscaling_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_appautoscaling_policy) ([#1798](https://github.com/turbot/steampipe-plugin-aws/pull/1798)) (Thanks [@jramosf](https://github.com/jramosf) for the contribution!)
  - [aws_identitystore_group_membership](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_identitystore_group_membership) ([#1782](https://github.com/turbot/steampipe-plugin-aws/pull/1782))
  - [aws_s3_bucket_intelligent_tiering_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_s3_bucket_intelligent_tiering_configuration) ([#1790](https://github.com/turbot/steampipe-plugin-aws/pull/1790))

_Enhancements_

- Added documentation on how to configure the plugin credentials when Steampipe is running on AWS ECS. Please refer [AssumeRole Credentials (in ECS)](https://hub.steampipe.io/plugins/turbot/aws#assumerole-credentials-in-ecs) for more information. ([#1800](https://github.com/turbot/steampipe-plugin-aws/pull/1800)) (Thanks [@Wade9320](https://github.com/Wade9320) for the contribution!)
- Added column `user_data` to `aws_ec2_launch_template_version` table. ([#1792](https://github.com/turbot/steampipe-plugin-aws/pull/1792))
- Added column `managed_actions` to `aws_elastic_beanstalk_environment` table. ([#1620](https://github.com/turbot/steampipe-plugin-aws/pull/1620))

_Bug fixes_

- Fixed `aws_acm_certificate` table to return certificates of all types of key algorithms instead of only the default `RSA_2048` algorithm. ([#1797](https://github.com/turbot/steampipe-plugin-aws/pull/1797))

## v0.107.0 [2023-06-21]

_What's new?_

- New tables added
  - [aws_ec2_managed_prefix_list_entry](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_managed_prefix_list_entry) ([#1781](https://github.com/turbot/steampipe-plugin-aws/pull/1781))
  - [aws_organizations_policy_target](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_policy_target) ([#1783](https://github.com/turbot/steampipe-plugin-aws/pull/1783))
  - [aws_route53_resolver_query_log_config](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_resolver_query_log_config) ([#1780](https://github.com/turbot/steampipe-plugin-aws/pull/1780))

_Enhancements_

- Added column `image_uri` to `aws_ecr_image` table. ([#1785](https://github.com/turbot/steampipe-plugin-aws/pull/1785))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#1775](https://github.com/turbot/steampipe-plugin-aws/pull/1775))

## v0.106.0 [2023-06-08]

_What's new?_

- New tables added
  - [aws_api_gateway_domain_name](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_api_gateway_domain_name) ([#1665](https://github.com/turbot/steampipe-plugin-aws/pull/1665))
  - [aws_route53_query_log](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_query_log) ([#1770](https://github.com/turbot/steampipe-plugin-aws/pull/1770))

_Bug fixes_

- Fixed the `ListConfig` of `aws_cloudformation_stack_resource` table to correctly return results instead of an empty row. ([#1771](https://github.com/turbot/steampipe-plugin-aws/pull/1771))

## v0.105.1 [2023-06-02]

_Bug fixes_

- Fixed the `associated_resources` column of `aws_wafv2_web_acl` table to also return associated CloudFront distributions. ([#1763](https://github.com/turbot/steampipe-plugin-aws/pull/1763))
- Fixed the syntax error in the example query of the `aws_inspector2_finding` table. ([#1764](https://github.com/turbot/steampipe-plugin-aws/pull/1764))

## v0.105.0 [2023-06-01]

_What's new?_

- New tables added
  - [aws_docdb_cluster_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_docdb_cluster_instance) ([#1755](https://github.com/turbot/steampipe-plugin-aws/pull/1755))
  - [aws_service_discovery_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_service_discovery_instance) ([#1741](https://github.com/turbot/steampipe-plugin-aws/pull/1741))
  - [aws_ssm_inventory_entry](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_inventory_entry) ([#1743](https://github.com/turbot/steampipe-plugin-aws/pull/1743))
  - [aws_sts_caller_identity](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sts_caller_identity) ([#1746](https://github.com/turbot/steampipe-plugin-aws/pull/1746))

_Bug fixes_

- Fixed `aws_inspector2_*` tables to correctly return data for all supported regions instead of only the `us-east-1` region. ([#1758](https://github.com/turbot/steampipe-plugin-aws/pull/1758))
- Fixed the `associated_resources` column in the `aws_wafv2_web_acl` table to include the associated resources of `API Gateway`, `App Sync`, and `Cognito User Pool`, in addition to the previously returned `Application Load Balancer resource type`. ([#1754](https://github.com/turbot/steampipe-plugin-aws/pull/1754))
- Fixed the `aws_wafv2_web_acl` table to return the missing `CloudFront` level web ACLs. ([#1752](https://github.com/turbot/steampipe-plugin-aws/pull/1752))

## v0.104.0 [2023-05-26]

_What's new?_

- New tables added
  - [aws_cloudwatch_metric_data_point](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_metric_data_point) ([#1655](https://github.com/turbot/steampipe-plugin-aws/pull/1655))
  - [aws_cloudwatch_metric_statistic_data_point](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_metric_statistic_data_point) ([#1649](https://github.com/turbot/steampipe-plugin-aws/pull/1649))
  - [aws_inspector2_coverage](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector2_coverage) ([#1657](https://github.com/turbot/steampipe-plugin-aws/pull/1657)) (Thanks [@jaredreisinger-drizly](https://github.com/jaredreisinger-drizly) for the contribution!!)
  - [aws_inspector2_coverage_statistics](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector2_coverage_statistics) ([#1657](https://github.com/turbot/steampipe-plugin-aws/pull/1657)) (Thanks [@jaredreisinger-drizly](https://github.com/jaredreisinger-drizly) for the contribution!!)
  - [aws_inspector2_finding](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector2_finding) ([#1657](https://github.com/turbot/steampipe-plugin-aws/pull/1657)) (Thanks [@jaredreisinger-drizly](https://github.com/jaredreisinger-drizly) for the contribution!!)
  - [aws_inspector2_member](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector2_member) ([#1657](https://github.com/turbot/steampipe-plugin-aws/pull/1657)) (Thanks [@jaredreisinger-drizly](https://github.com/jaredreisinger-drizly) for the contribution!!)
  - [aws_rds_db_instance_automated_backup](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_automated_backup) ([#1721](https://github.com/turbot/steampipe-plugin-aws/pull/1721))

_Enhancements_

- Added an example query in aws_iam_role table doc. ([#1745](https://github.com/turbot/steampipe-plugin-aws/pull/1745))

## v0.103.0 [2023-05-18]

_What's new?_

- New tables added
  - [aws_ec2_client_vpn_endpoint](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_client_vpn_endpoint) ([#1722](https://github.com/turbot/steampipe-plugin-aws/pull/1722))
  - [aws_ec2_launch_template_version](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_launch_template_version) ([#1725](https://github.com/turbot/steampipe-plugin-aws/pull/1725))
  - [aws_service_discovery_namespace](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_service_discovery_namespace) ([#1735](https://github.com/turbot/steampipe-plugin-aws/pull/1735))
  - [aws_service_discovery_service](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_service_discovery_service) ([#1739](https://github.com/turbot/steampipe-plugin-aws/pull/1739))
  - [aws_servicecatalog_product](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicecatalog_product) ([#1638](https://github.com/turbot/steampipe-plugin-aws/pull/1638))
  - [aws_ssm_managed_instance_patch_state](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_managed_instance_patch_state) ([#1732](https://github.com/turbot/steampipe-plugin-aws/pull/1732))
  - [aws_wellarchitected_answer](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_answer) ([#1699](https://github.com/turbot/steampipe-plugin-aws/pull/1699))
  - [aws_wellarchitected_check_detail](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_check_detail) ([#1700](https://github.com/turbot/steampipe-plugin-aws/pull/1700))
  - [aws_wellarchitected_check_summary](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_check_summary) ([#1700](https://github.com/turbot/steampipe-plugin-aws/pull/1700))
  - [aws_wellarchitected_consolidated_report](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_consolidated_report) ([#1704](https://github.com/turbot/steampipe-plugin-aws/pull/1704))
  - [aws_wellarchitected_lens_review_improvement](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_lens_review_improvement) ([#1695](https://github.com/turbot/steampipe-plugin-aws/pull/1695))
  - [aws_wellarchitected_lens_review_report](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_lens_review_report) ([#1697](https://github.com/turbot/steampipe-plugin-aws/pull/1697))
  - [aws_wellarchitected_lens_share](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_lens_share) ([#1698](https://github.com/turbot/steampipe-plugin-aws/pull/1698))
  - [aws_wellarchitected_share_invitation](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_share_invitation) ([#1692](https://github.com/turbot/steampipe-plugin-aws/pull/1692))

_Bug fixes_

- Fixed the `source_account_id` optional qual column definition in `aws_security_hub_finding` table. ([#1737](https://github.com/turbot/steampipe-plugin-aws/pull/1737)) (Thanks [@gabrielsoltz](https://github.com/gabrielsoltz) for the contribution!)
- Fixed the example query in the doc for the `aws_ssoadmin_account_assignment` table. ([#1734](https://github.com/turbot/steampipe-plugin-aws/pull/1734))

## v0.102.0 [2023-05-11]

_What's new?_

- New tables added
  - [aws_config_retention_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_config_retention_configuration) ([#1718](https://github.com/turbot/steampipe-plugin-aws/pull/1718))

_Enhancements_

- Added column `repository_scanning_configuration` to `aws_ecr_repository` table. ([#1719](https://github.com/turbot/steampipe-plugin-aws/pull/1719))
- Added column `source_account_id` to `aws_securityhub_finding` table. ([#1703](https://github.com/turbot/steampipe-plugin-aws/pull/1703)) (Thanks [@gabrielsoltz](https://github.com/gabrielsoltz) for the contribution!)

_Bug fixes_

- Fixed `aws_ecr_image_scan_finding` table to return an empty row instead of an error when image scanning is in progress. ([#1728](https://github.com/turbot/steampipe-plugin-aws/pull/1728)) (Thanks [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for the contribution!)
- Fixed the `GetConfig` of the `aws_ssm_document` table to use `arn` instead of `name` as a key column to avoid failures in querying multiple regions with the same document name. ([#1720](https://github.com/turbot/steampipe-plugin-aws/pull/1720))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#1685](https://github.com/turbot/steampipe-plugin-aws/pull/1685))

## v0.101.0 [2023-04-25]

_What's new?_

- New tables added
  - [aws_wellarchitected_lens](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_lens) ([#1689](https://github.com/turbot/steampipe-plugin-aws/pull/1689))
  - [aws_wellarchitected_notification](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_notification) ([#1693](https://github.com/turbot/steampipe-plugin-aws/pull/1693))
  - [aws_wellarchitected_workload_share](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_workload_share) ([#1690](https://github.com/turbot/steampipe-plugin-aws/pull/1690))

_Enhancements_

- Added `maintenance_options`, `licenses`, `placement_affinity`, `placement_group_id`, `placement_host_id`, `placement_host_resource_group_arn`, `placement_partition_number`, and `spot_instance_request_id` columns to `aws_ec2_instance` table. ([#1709](https://github.com/turbot/steampipe-plugin-aws/pull/1709))
- Added `workspace` column to `aws_wellarchitected_milestone` table.
- Removed hydrate requirement for `milestone_number` column in `aws_wellarchitected_lens_review` table.

## v0.100.0 [2023-04-15]

_What's new?_

- New tables added
  - [aws_s3_object](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_s3_object) ([#1623](https://github.com/turbot/steampipe-plugin-aws/pull/1623))
  - [aws_wellarchitected_lens_review](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_lens_review) ([#1694](https://github.com/turbot/steampipe-plugin-aws/pull/1694))
  - [aws_wellarchitected_milestone](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_milestone) ([#1696](https://github.com/turbot/steampipe-plugin-aws/pull/1696))

## v0.99.0 [2023-04-07]

_What's new?_

- New tables added
  - [aws_ssoadmin_account_assignment](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssoadmin_account_assignment) ([#1673](https://github.com/turbot/steampipe-plugin-aws/pull/1673)) (Thanks [@janslow](https://github.com/janslow) for the contribution!)
  - [aws_athena_query_execution](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_athena_query_execution) ([#1666](https://github.com/turbot/steampipe-plugin-aws/pull/1666)) (Thanks [@rinzool](https://github.com/rinzool) for the contribution!)
  - [aws_athena_workgroup](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_athena_workgroup) ([#1666](https://github.com/turbot/steampipe-plugin-aws/pull/1666)) (Thanks [@rinzool](https://github.com/rinzool) for the contribution!)

_Bug fixes_

- Fixed typos in the `ListConfig` of `aws_sfn_state_machine_*` tables. ([#1686](https://github.com/turbot/steampipe-plugin-aws/pull/1686)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Fixed the data type of `tags` column of `aws_securitylake_data_lake` and `aws_simspaceweaver_simulation` tables to be of `JSON` type instead of `STRING`. ([#1683](https://github.com/turbot/steampipe-plugin-aws/pull/1683))
- Fixed the `aws_organizations_policy` table to correctly return all the organization policies instead of duplicate data. ([#1681](https://github.com/turbot/steampipe-plugin-aws/pull/1681))

## v0.98.0 [2023-03-31]

_What's new?_

- New tables added
  - [aws_codedeploy_deployment_config](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codedeploy_deployment_config) ([#1662](https://github.com/turbot/steampipe-plugin-aws/pull/1662))
  - [aws_organizations_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_policy) ([#1641](https://github.com/turbot/steampipe-plugin-aws/pull/1641))
  - [aws_s3_multi_region_access_point](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_s3_multi_region_access_point) ([#1486](https://github.com/turbot/steampipe-plugin-aws/pull/1486))
  - [aws_ssm_document_permission](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_document_permission) ([#1670](https://github.com/turbot/steampipe-plugin-aws/pull/1670))
  - [aws_vpc_eip_address_transfer](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_eip_address_transfer) ([#1521](https://github.com/turbot/steampipe-plugin-aws/pull/1521))
  - [aws_wafregional_rule_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafregional_rule_group) ([#1661](https://github.com/turbot/steampipe-plugin-aws/pull/1661))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#1676](https://github.com/turbot/steampipe-plugin-aws/pull/1676))

## v0.97.0 [2023-03-24]

_What's new?_

- New tables added
  - [aws_appstream_image](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_appstream_image) ([#1663](https://github.com/turbot/steampipe-plugin-aws/pull/1663))
  - [aws_codedeploy_deployment_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codedeploy_deployment_group) ([#1658](https://github.com/turbot/steampipe-plugin-aws/pull/1658))
  - [aws_cost_by_tag](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_tag) ([#1536](https://github.com/turbot/steampipe-plugin-aws/pull/1536))
  - [aws_networkfirewall_firewall](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_networkfirewall_firewall) ([#1630](https://github.com/turbot/steampipe-plugin-aws/pull/1630))
  - [aws_wafregional_web_acl](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafregional_web_acl) ([#1660](https://github.com/turbot/steampipe-plugin-aws/pull/1660))

_Bug fixes_

- Fixed the `aws_health_affected_entity` table to correctly return results instead of an error. ([#1659](https://github.com/turbot/steampipe-plugin-aws/pull/1659))

## v0.96.0 [2023-03-10]

_What's new?_

- New tables added
  - [aws_cloudformation_stack_resource](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudformation_stack_resource) ([#1634](https://github.com/turbot/steampipe-plugin-aws/pull/1634))

_Enhancements_

- Added columns `dkim_attributes` and `identity_mail_from_domain_attributes` to `aws_ses_domain_identity` table. ([#1640](https://github.com/turbot/steampipe-plugin-aws/pull/1640))

_Bug fixes_

- Fixed `aws_cloudfront_response_headers_policy` table to remove duplicate results. ([#1642](https://github.com/turbot/steampipe-plugin-aws/pull/1642)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v520-2023-03-02) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#1609](https://github.com/turbot/steampipe-plugin-aws/pull/1609))

## v0.95.0 [2023-03-03]

_What's new?_

- New tables added
  - [aws_codebuild_build](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codebuild_build) ([#1608](https://github.com/turbot/steampipe-plugin-aws/pull/1608))
  - [aws_emr_block_public_access_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_block_public_access_configuration) ([#1602](https://github.com/turbot/steampipe-plugin-aws/pull/1602))
  - [aws_servicecatalog_portfolio](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicecatalog_portfolio) ([#1636](https://github.com/turbot/steampipe-plugin-aws/pull/1636))

_Bug fixes_

- Fixed the `aws_cloudfront_function` table to correctly return data instead of an error when a `name` is passed in the `where` clause. ([#1628](https://github.com/turbot/steampipe-plugin-aws/pull/1628))
- Fixed the `aws_guardduty_ipset` table to correctly return all the IPsets instead of a panic interface conversion error. ([#1627](https://github.com/turbot/steampipe-plugin-aws/pull/1627))
- Fixed the API limits of the `aws_glue_security_configuration` table to correctly return data instead of an error. ([#1626](https://github.com/turbot/steampipe-plugin-aws/pull/1626))

## v0.94.0 [2023-02-25]

_What's new?_

- New tables added
  - [aws_ec2_launch_template](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_launch_template) ([#1543](https://github.com/turbot/steampipe-plugin-aws/pull/1543))

_Enhancements_

- Added column `data_protection` and `data_protection_policy` to `aws_cloudwatch_log_group` table. ([#1483](https://github.com/turbot/steampipe-plugin-aws/pull/1483))
- Added column `website_configuration` to `aws_s3_bucket` table. ([#1618](https://github.com/turbot/steampipe-plugin-aws/pull/1618))
- Added column `object_ownership_controls` to `aws_s3_bucket` table. ([#1548](https://github.com/turbot/steampipe-plugin-aws/pull/1548))
- Added column `launch_template_data` to `aws_ec2_instance` table. ([#1553](https://github.com/turbot/steampipe-plugin-aws/pull/1553))
- Added column `tracing_config` to `aws_lambda_function` table. ([#1601](https://github.com/turbot/steampipe-plugin-aws/pull/1601))
- Updated Parliament IAM permissions to the latest. ([#1599](https://github.com/turbot/steampipe-plugin-aws/pull/1599))

_Bug fixes_

- Fixed the `title` column in `aws_api_gatewayv2_route` table to correctly return data instead of `null`. ([#1568](https://github.com/turbot/steampipe-plugin-aws/pull/1568))
- Fixed the `tags_src` column in `aws_cloudformation_stack` table to correctly return raw tag data instead of a formatted one. ([#1568](https://github.com/turbot/steampipe-plugin-aws/pull/1568))
- Fixed the `architectures`, `file_system_configs` and `snap_start` columns in `aws_lambda_function` table to correctly return data instead of `null`. ([#1619](https://github.com/turbot/steampipe-plugin-aws/pull/1619))
- Fixed `aws_ec2_managed_prefix_list` table to return an empty row instead of an error in unsupported `me-south-1` region. ([#1577](https://github.com/turbot/steampipe-plugin-aws/pull/1577))
- Fixed the `aws_eventbridge_rule` table to return rules for all the event bridges instead of only default event bridges. ([#1590](https://github.com/turbot/steampipe-plugin-aws/pull/1590)) (Thanks [@brentmitchell25](https://github.com/brentmitchell25) for the fix!!)

## v0.93.0 [2023-02-17]

_What's new?_

- Added `default_region` config arg, which allows you to set your preferred (closest) region to optimize API calls to global resources. ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))

_Enhancements_

- EC2 Role & SSO credentials are now used until they expire, reducing throttling & reloading. ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))
- Optimized API calls to use the default region, reducing latency for common APIs. ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))
- Optimized caching to reduce race conditions & extend timeouts (e.g. credentials). ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))
- Optimized per-region API calls to regions supported by the service only. ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))
- Optimized API client management to one per account, instead of one per region. ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0-rc.7](https://github.com/turbot/steampipe-plugin-sdk/blob/e26550be74f53fa902dbc66576d3f828980556c3/CHANGELOG.md#v520-tbd) which includes additional cache function wrappers and matrix function improvements. ([#1559](https://github.com/turbot/steampipe-plugin-aws/pull/1559))

## v0.92.2 [2023-02-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.12](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4112-2023-02-09) which fixes the query caching functionality. ([#1578](https://github.com/turbot/steampipe-plugin-aws/pull/1578))

## v0.92.1 [2023-01-24]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.11](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4111-2023-01-24) which fixes the issue of non-caching of all the columns of the queried table. ([#1557](https://github.com/turbot/steampipe-plugin-aws/pull/1557))

## v0.92.0 [2023-01-19]

_What's new?_

- New tables added
  - [aws_api_gatewayv2_route](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_api_gatewayv2_route) ([#1544](https://github.com/turbot/steampipe-plugin-aws/pull/1544))
  - [aws_health_affected_entity](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_health_affected_entity) ([#1525](https://github.com/turbot/steampipe-plugin-aws/pull/1525))

_Enhancements_

- Added column `access_log_settings` to `aws_api_gatewayv2_stage` table. ([#1546](https://github.com/turbot/steampipe-plugin-aws/pull/1546))

_Bug fixes_

- Fixed the `aws_ec2_ami` table to only return images owned by the AWS account. ([#1535](https://github.com/turbot/steampipe-plugin-aws/pull/1535))
- Fixed the `aws_ec2_ami_shared` table to return images from any AWS account (images owned by the AWS account or shared by other accounts) when either an `owner_id` or an `image_id` or both the parameters are passed in the `where` clause. ([#1535](https://github.com/turbot/steampipe-plugin-aws/pull/1535))

## v0.91.1 [2023-01-17]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.9](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v419-2022-11-30) which fixes hydrate function caching for aggregator connections. ([#1540](https://github.com/turbot/steampipe-plugin-aws/pull/1540))

## v0.91.0 [2023-01-09]

_What's new?_

- New tables added
  - [aws_glue_data_quality_ruleset](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_data_quality_ruleset) ([#1513](https://github.com/turbot/steampipe-plugin-aws/pull/1513))

_Bug fixes_

- Fixed `aws_s3_access_point` table to return access points from all the configured regions instead of only `us-east-1`. ([#1522](https://github.com/turbot/steampipe-plugin-aws/pull/1522))
- Fixed the `aws_ebs_snapshot` table to return snapshots from different AWS accounts when an `owner_alias` or an `owner_id` or a `snapshot_id` is passed in the `where` clause. ([#1530](https://github.com/turbot/steampipe-plugin-aws/pull/1530))

## v0.90.0 [2022-12-28]

_What's new?_

- New tables added
  - [aws_cloudtrail_channel](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_channel) ([#1517](https://github.com/turbot/steampipe-plugin-aws/pull/1517))
  - [aws_cloudtrail_import](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_import) ([#1515](https://github.com/turbot/steampipe-plugin-aws/pull/1515))
  - [aws_securitylake_data_lake](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securitylake_data_lake) ([#1520](https://github.com/turbot/steampipe-plugin-aws/pull/1520))

_Bug fixes_

- Fixed the `aws_api_gatewayv2_*` tables to return an empty row for unsupported region `ap-southeast-3` instead of an error. ([#1527](https://github.com/turbot/steampipe-plugin-aws/pull/1527))

## v0.89.0 [2022-12-23]

_What's new?_

- New tables added
  - [aws_cloudtrail_query](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_query) ([#1437](https://github.com/turbot/steampipe-plugin-aws/pull/1437))
  - [aws_drs_job](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_drs_job) ([#1492](https://github.com/turbot/steampipe-plugin-aws/pull/1492))
  - [aws_drs_recovery_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_drs_recovery_instance) ([#1493](https://github.com/turbot/steampipe-plugin-aws/pull/1493))
  - [aws_drs_recovery_snapshot](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_drs_recovery_snapshot) ([#1494](https://github.com/turbot/steampipe-plugin-aws/pull/1494))
  - [aws_mgn_application](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_mgn_application) ([#1499](https://github.com/turbot/steampipe-plugin-aws/pull/1499))
  - [aws_oam_link](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_oam_link) ([#1498](https://github.com/turbot/steampipe-plugin-aws/pull/1498))
  - [aws_oam_sink](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_oam_sink) ([#1495](https://github.com/turbot/steampipe-plugin-aws/pull/1495))
  - [aws_simspaceweaver_simulation](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_simspaceweaver_simulation) ([#1501](https://github.com/turbot/steampipe-plugin-aws/pull/1501))

_Enhancements_

- Added column `addon_configuration` to `aws_eks_addon_version` table. ([#1514](https://github.com/turbot/steampipe-plugin-aws/pull/1514))
- Added column `standards_managed_by` to `aws_securityhub_standards_subscription` table. ([#1511](https://github.com/turbot/steampipe-plugin-aws/pull/1511))
- Added column `launch_configuration` to `aws_drs_source_server` table. ([#1496](https://github.com/turbot/steampipe-plugin-aws/pull/1496))
- Added column `protection` to `aws_ecs_task` table. ([#1500](https://github.com/turbot/steampipe-plugin-aws/pull/1500))

_Bug fixes_

- Fixed the `insight_selectors` column in `aws_cloudtrail_trail` table to correctly return data instead of `nil`. ([#1512](https://github.com/turbot/steampipe-plugin-aws/pull/1512))
- Fixed the `tags` and `tags_src` column in `aws_dynamodb_table` table to correctly handle the `ResourceNotFoundException` error and return `nil` when an invalid `arn` is passed in the where clause. ([#1518](https://github.com/turbot/steampipe-plugin-aws/pull/1518))

## v0.88.0 [2022-12-15]

_What's new?_

- New tables added
  - [aws_cloudtrail_event_data_store](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_event_data_store) ([#1433](https://github.com/turbot/steampipe-plugin-aws/pull/1433))
  - [aws_drs_source_server](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_drs_source_server) ([#1465](https://github.com/turbot/steampipe-plugin-aws/pull/1465))
  - [aws_pipes_pipe](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_pipes_pipe) ([#1487](https://github.com/turbot/steampipe-plugin-aws/pull/1487))
  - [aws_vpc_verified_access_endpoint](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_verified_access_endpoint) ([#1491](https://github.com/turbot/steampipe-plugin-aws/pull/1491))
  - [aws_vpc_verified_access_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_verified_access_group) ([#1485](https://github.com/turbot/steampipe-plugin-aws/pull/1485))
  - [aws_vpc_verified_access_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_verified_access_instance) ([#1484](https://github.com/turbot/steampipe-plugin-aws/pull/1484))
  - [aws_vpc_verified_access_trust_provider](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_verified_access_trust_provider) ([#1482](https://github.com/turbot/steampipe-plugin-aws/pull/1482))

_Enhancements_

- Added column `platform_family` to `aws_ecs_service` table. ([#1490](https://github.com/turbot/steampipe-plugin-aws/pull/1490))

## v0.87.0 [2022-12-02]

_Breaking changes_

- The `aws_cloudwatch_metric` table rows now contain a CloudWatch metric each, instead of a dimension name/value pair. Dimensions for each metric can be found in the `dimensions` column and to filter on specific dimensions, you can pass dimensions through the `dimensions_filter` key column. Please see [aws_cloudwatch_metric Examples](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_metric#examples) for query examples using the new columns.
- Renamed column `name` to `metric_name` in the `aws_cloudwatch_metric` table.

_What's new?_

- New tables added
  - [aws_backup_legal_hold](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_legal_hold) ([#1464](https://github.com/turbot/steampipe-plugin-aws/pull/1464))
  - [aws_dax_parameter](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dax_parameter) ([#1434](https://github.com/turbot/steampipe-plugin-aws/pull/1434))
  - [aws_securitylake_subscriber](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securitylake_subscriber) ([#1463](https://github.com/turbot/steampipe-plugin-aws/pull/1463))

_Enhancements_

- Added `evaluation_modes` column to the `aws_config_rule` table. ([#1476](https://github.com/turbot/steampipe-plugin-aws/pull/1476))
- Added `snap_start` column to the `aws_lambda_function` table. ([#1477](https://github.com/turbot/steampipe-plugin-aws/pull/1477))
- Added `capacity_allocations` column to the `aws_ec2_capacity_reservation` table. ([#1428](https://github.com/turbot/steampipe-plugin-aws/pull/1428))
- Added `imds_support` column to `aws_ec2_ami` and `aws_ec2_ami_shared` tables. ([#1430](https://github.com/turbot/steampipe-plugin-aws/pull/1430))

## v0.86.0 [2022-11-28]

_What's new?_

- New tables added
  - [aws_dax_parameter_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dax_parameter_group) ([#1426](https://github.com/turbot/steampipe-plugin-aws/pull/1426))

_Bug fixes_

- Fixed the `aws_rds_db_proxy table` table to return empty rows for unsupported regions instead of an error. ([#1427](https://github.com/turbot/steampipe-plugin-aws/pull/1427))

## v0.85.0 [2022-11-24]

_What's new?_

- New tables added
  - [aws_cloudsearch_domain](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudsearch_domain) ([#1413](https://github.com/turbot/steampipe-plugin-aws/pull/1413))
  - [aws_eks_fargate_profile](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_fargate_profile) ([#1409](https://github.com/turbot/steampipe-plugin-aws/pull/1409))
  - [aws_health_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_health_event) ([#1167](https://github.com/turbot/steampipe-plugin-aws/pull/1167))
  - [aws_kms_alias](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_kms_alias) ([#1405](https://github.com/turbot/steampipe-plugin-aws/pull/1405))
  - [aws_rds_db_proxy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_proxy) ([#1414](https://github.com/turbot/steampipe-plugin-aws/pull/1414))

_Bug fixes_

- Fixed the `ServiceNotFoundException` error in the `aws_ecs_task` table to return an empty row when an invalid value is passed in the `service_name` filter. ([#1418](https://github.com/turbot/steampipe-plugin-aws/pull/1418))
- Fixed the `ResourceNotFoundException` in the `aws_cloudwatch_log_metric_filter` table to return an empty row when an invalid value is passed in the `log_group_name` filter. ([#1420](https://github.com/turbot/steampipe-plugin-aws/pull/1420))

## v0.84.2 [2022-11-22]

_Bug fixes_

- Fixed the plugin to use environment variables like `AWS_REGION`, `AWS_DEFAULT_REGION` etc., when no regions are specified in the `aws.spc` file. ([#1411](https://github.com/turbot/steampipe-plugin-aws/pull/1411))

## v0.84.1 [2022-11-18]

_Dependencies_

- Recompiled plugin with [aws-sdk-go v1.44.141](https://github.com/aws/aws-sdk-go/blob/main/CHANGELOG.md#release-v144141-2022-11-18) and [aws-sdk-go-v2/service/route53 v1.24.0](https://github.com/aws/aws-sdk-go-v2/blob/main/service/route53/CHANGELOG.md#v1240-2022-11-15) to update service endpoints.

## v0.84.0 [2022-11-17]

_Enhancements_

- Improved default region checking for global and region limited services. ([#1397](https://github.com/turbot/steampipe-plugin-aws/pull/1397))

## v0.83.0 [2022-11-16]

_What's new?_

- New tables added
  - [aws_ec2_spot_price](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_spot_price) ([#1378](https://github.com/turbot/steampipe-plugin-aws/pull/1378)) (Thanks to [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for the new table!)
  - [aws_iam_service_specific_credential](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iam_service_specific_credential) ([#1390](https://github.com/turbot/steampipe-plugin-aws/pull/1390))
  - [aws_pricing_product](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_pricing_product) ([#1369](https://github.com/turbot/steampipe-plugin-aws/pull/1369)) (Thanks to [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for another new table!)
  - [aws_resource_explorer_index](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_index) ([#1396](https://github.com/turbot/steampipe-plugin-aws/pull/1396))
  - [aws_resource_explorer_search](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_search) ([#1396](https://github.com/turbot/steampipe-plugin-aws/pull/1396))
  - [aws_resource_explorer_supported_resource_type](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) ([#1396](https://github.com/turbot/steampipe-plugin-aws/pull/1396))

_Bug fixes_

- Fixed queries failing for `aws_s3_access_point` table when an invalid bucket name is specified. ([#1395](https://github.com/turbot/steampipe-plugin-aws/pull/1395))

## v0.82.0 [2022-11-09]

_Enhancements_

- Added `workflow_status` column to the `aws_securityhub_finding` table. ([#1377](https://github.com/turbot/steampipe-plugin-aws/pull/1377)) (Thanks [@gabrielsoltz](https://github.com/gabrielsoltz) for the contribution!)

_Bug fixes_

- Fixed the `aws_api_gatewayv2_*` tables to correctly return results instead of an error by skipping the unsupported `me-central-1` region. ([#1388](https://github.com/turbot/steampipe-plugin-aws/pull/1388))
- Fixed the `billing_mode` column in `aws_dynamodb_table` to correctly return results instead of an error. ([#1387](https://github.com/turbot/steampipe-plugin-aws/pull/1387))

_Deprecated_

- Deprecated the `workflow_state` column in the `aws_securityhub_finding` table per [AWS documentation](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/securityhub/get-findings.html#options). We recommend updating any workflows and queries to use `workflow_status` instead of `workflow_state`. ([#1377](https://github.com/turbot/steampipe-plugin-aws/pull/1377))

## v0.81.1 [2022-11-09]

_Bug fixes_

- Fixed the typo in the example query of `aws_efs_file_system` table document to use `ValueInStandard` instead of `ValueInIA`. ([#1381](https://github.com/turbot/steampipe-plugin-aws/pull/1381)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#1382](https://github.com/turbot/steampipe-plugin-aws/pull/1382))

## v0.81.0 [2022-11-04]

_Enhancements_

- Added `set_identifier` as an optional list key column in `aws_route53_record` table. ([#1375](https://github.com/turbot/steampipe-plugin-aws/pull/1375))
- Updated 30+ tables to use [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2). ([#1361](https://github.com/turbot/steampipe-plugin-aws/pull/1361)) ([#1371](https://github.com/turbot/steampipe-plugin-aws/pull/1371))

_Bug fixes_

- Fixed paging in `aws_route53_record` table to ensure all records are returned. ([#1375](https://github.com/turbot/steampipe-plugin-aws/pull/1375))
- Fixed invalid pointer usage causing duplicate values in `attribute_name` column for `aws_pricing_service_attribute` table. ([#1372](https://github.com/turbot/steampipe-plugin-aws/pull/1372)) (Thanks to [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for the fix!)
- Fixed example queries in `aws_ebs_volume` table document. ([#1368](https://github.com/turbot/steampipe-plugin-aws/pull/1368))

## v0.80.0 [2022-10-21]

_What's new?_

- New tables added
  - [aws_ecr_image_scan_finding](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecr_image_scan_finding) ([#1315](https://github.com/turbot/steampipe-plugin-aws/pull/1315)) (Thanks to [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for the new table!)
  - [aws_lightsail_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lightsail_instance) ([#1359](https://github.com/turbot/steampipe-plugin-aws/pull/1359))

_Enhancements_

- Added `owner_type` column to the `aws_ssm_document` table to allow filtering on SSM documents by AWS account type. ([#1337](https://github.com/turbot/steampipe-plugin-aws/pull/1337))
- Updated 80+ tables to use [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2). ([#1337](https://github.com/turbot/steampipe-plugin-aws/pull/1337)) ([#1357](https://github.com/turbot/steampipe-plugin-aws/pull/1357))

_Bug fixes_

- Fixed `status` column type from JSON to string in `aws_ssm_association` table. ([#1337](https://github.com/turbot/steampipe-plugin-aws/pull/1337))
- Removed unsupported `TAGS` dimension note in `aws_cost_usage` table doc. ([#1362](https://github.com/turbot/steampipe-plugin-aws/pull/1362))

_Deprecated_

- Deprecated `image_details` and `image_scanning_findings` columns in `aws_ecr_repository` table to avoid throttling issues. Please use the `aws_ecr_image` and `aws_ecr_image_scan_finding` tables instead. ([#1198](https://github.com/turbot/steampipe-plugin-aws/pull/1198))

## v0.79.1 [2022-10-17]

_Bug fixes_

- Fixed unsupported region check in `aws_dlm_lifecycle_policy` table to allow queries for valid regions.
- Fixed paging in `aws_route53_record` table to return all records correctly. ([#1356](https://github.com/turbot/steampipe-plugin-aws/pull/1356))

## v0.79.0 [2022-10-14]

_Enhancements_

- Updated 70+ tables to use [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2). ([#1324](https://github.com/turbot/steampipe-plugin-aws/pull/1324))
- Added `managed_actions` column to `aws_elastic_beanstalk_environment` table. ([#996](https://github.com/turbot/steampipe-plugin-aws/pull/996))
- Added the following columns to the `aws_ec2_instance` table:
  - `ami_launch_index`
  - `architecture`
  - `boot_mode`
  - `capacity_reservation_id`
  - `capacity_reservation_specification`
  - `client_token`
  - `ena_support`
  - `enclave_options`
  - `hibernation_options`
  - `platform`
  - `platform_details`
  - `private_dns_name_options`
  - `state_transition_reason`
  - `tpm_support`
  - `usage_operation`
  - `usage_operation_update_time`

_Bug fixes_

- Removed duplicate values in `inline_policies` column in `aws_iam_role` and `aws_iam_user` tables. ([#1346](https://github.com/turbot/steampipe-plugin-aws/pull/1346))
- Fixed queries failing for the `aws_acm_certificate` table when querying the `title` column. ([#1351](https://github.com/turbot/steampipe-plugin-aws/pull/1351))
- Fixed empty check for `regions` config arg incorrectly failing when at least 1 other config arg is set. ([#1349](https://github.com/turbot/steampipe-plugin-aws/pull/1349))
- Fixed queries that specify `service_name` for the `aws_ecs_task` table returning no rows if an unqualified query was run first. ([#1338](https://github.com/turbot/steampipe-plugin-aws/pull/1338))

## v0.78.0 [2022-09-23]

_What's new?_

- New tables added
  - [aws_account_alternate_contact](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_account_alternate_contact) ([#1310](https://github.com/turbot/steampipe-plugin-aws/pull/1310))
  - [aws_account_contact](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_account_contact) ([#1310](https://github.com/turbot/steampipe-plugin-aws/pull/1310))
  - [aws_msk_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_msk_cluster) ([#1291](https://github.com/turbot/steampipe-plugin-aws/pull/1291))
  - [aws_msk_serverless_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_msk_serverless_cluster) ([#1291](https://github.com/turbot/steampipe-plugin-aws/pull/1291))

_Enhancements_

- Updated index doc **Configuring AWS Credentials** section to use consistent profile and account names. ([#1209](https://github.com/turbot/steampipe-plugin-aws/pull/1209)) (Thanks to [@michael-ullrich-1010](https://github.com/michael-ullrich-1010) for the contribution!)
- Improved plugin error message when the `regions` config argument is set to an invalid value `[]`.

_Bug fixes_

- `aws_macie2_classification_job` table now checks for supported regions.

## v0.77.0 [2022-09-15]

_What's new?_

- New tables added
  - [aws_appconfig_application](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_appconfig_application) ([#1253](https://github.com/turbot/steampipe-plugin-aws/pull/1253))
  - [aws_codeartifact_domain](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codeartifact_domain) ([#1184](https://github.com/turbot/steampipe-plugin-aws/pull/1184))
  - [aws_codeartifact_repository](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codeartifact_repository) ([#1268](https://github.com/turbot/steampipe-plugin-aws/pull/1268))
  - [aws_codedeploy_app](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codedeploy_app) ([#1295](https://github.com/turbot/steampipe-plugin-aws/pull/1295))
  - [aws_redshiftserverless_namespace](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshiftserverless_namespace) ([#1305](https://github.com/turbot/steampipe-plugin-aws/pull/1305))
  - [aws_redshiftserverless_workgroup](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshiftserverless_workgroup) ([#1304](https://github.com/turbot/steampipe-plugin-aws/pull/1304))

_Enhancements_

- Added `access_key_last_used_date`, `access_key_last_used_region` columns and `access_key_last_used_service` to `aws_iam_access_key` table. ([#1281](https://github.com/turbot/steampipe-plugin-aws/pull/1281))
- Added `vpc_endpoint_connections` column to `aws_vpc_endpoint_service` table. ([#1104](https://github.com/turbot/steampipe-plugin-aws/pull/1104))
- Updated the following tables to use [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2):
  - `aws_vpc_route_table`

_Bug fixes_

- `aws_dynamodb_table_export` table queries no longer fail when passing in `arn` get key column.
- `aws_ec2_transit_gateway`, `aws_ec2_transit_gateway_route`, `aws_ec2_transit_gateway_route_table`, and `aws_ec2_transit_gateway_vpc_attachment` tables should not error in me-central-1 region. ([#1282](https://github.com/turbot/steampipe-plugin-aws/pull/1282))
- `aws_vpc_eip` table now handles EIPs in EC2-Classic properly. ([#1308](https://github.com/turbot/steampipe-plugin-aws/pull/1308))
- `aws_wafregional_rule` table now properly checks for supported regions. ([#1306](https://github.com/turbot/steampipe-plugin-aws/pull/1306))

_Deprecated_

- Deprecated `verification_token` column in `aws_ses_email_identity` table since there is no verification token for email identities. This column will be removed in a future version.

## v0.76.0 [2022-09-09]

_What's new?_

- New tables added
  - [aws_cloudwatch_log_subscription_filter](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_subscription_filter) ([#1243](https://github.com/turbot/steampipe-plugin-aws/pull/1243))
  - [aws_dax_subnet_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dax_subnet_group) ([#1298](https://github.com/turbot/steampipe-plugin-aws/pull/1298))
  - [aws_docdb_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_docdb_cluster) ([#1019](https://github.com/turbot/steampipe-plugin-aws/pull/1019))
  - [aws_globalaccelerator_accelerator](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_globalaccelerator_accelerator) ([#1091](https://github.com/turbot/steampipe-plugin-aws/pull/1091)) (Thanks to [@nmische](https://github.com/nmische) for the contribution!)
  - [aws_globalaccelerator_endpoint_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_globalaccelerator_endpoint_group) ([#1091](https://github.com/turbot/steampipe-plugin-aws/pull/1091))
  - [aws_globalaccelerator_listener](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_globalaccelerator_listener) ([#1091](https://github.com/turbot/steampipe-plugin-aws/pull/1091))

_Enhancements_

- Added column `code` to `aws_lambda_function` table. ([#1293](https://github.com/turbot/steampipe-plugin-aws/pull/1293))
- Updated the `title` column of `aws_kms_key` table to first use the key alias if available, else fall back to the key ID. ([#1246](https://github.com/turbot/steampipe-plugin-aws/pull/1246))

_Bug fixes_

- Fixed the `url_config` column in `aws_lambda_function` table to return `null` instead of an access denied exception errors for US Government cloud regions. ([#1285](https://github.com/turbot/steampipe-plugin-aws/pull/1285))
- Fixed the `sns_topic_arn` column in `aws_backup_vault` table to correctly return a value instead of `null`. ([#1280](https://github.com/turbot/steampipe-plugin-aws/pull/1280))
- Fixed all the tables of CodeBuild and Serverless Application Repository services to return empty rows instead of an error for unsupported regions. ([#1289](https://github.com/turbot/steampipe-plugin-aws/pull/1289))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which fixes incorrect cache hits in multi-region queries which use the `region` column in the where clause. ([#387](https://github.com/turbot/steampipe-plugin-gcp/pull/387))

## v0.75.1 [2022-08-31]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.5](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v415-2022-08-31) which includes connection cache TTL fixes.

## v0.75.0 [2022-08-30]

_What's new?_

- New tables added
  - [aws_ecr_image](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecr_image) ([#1200](https://github.com/turbot/steampipe-plugin-aws/pull/1200))

_Enhancements_

- Added column `disable_execute_api_endpoint` to `aws_api_gatewayv2_api` table. ([#1242](https://github.com/turbot/steampipe-plugin-aws/pull/1242))
- Updated the following tables to use [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2): ([#1219](https://github.com/turbot/steampipe-plugin-aws/pull/1219))
  - `aws_cost_by_account_daily`
  - `aws_cost_by_account_monthly`
  - `aws_cost_by_record_type_daily`
  - `aws_cost_by_record_type_monthly`
  - `aws_cost_by_service_daily`
  - `aws_cost_by_service_monthly`
  - `aws_cost_by_service_usage_type_daily`
  - `aws_cost_by_service_usage_type_monthly`
  - `aws_cost_forecast_daily`
  - `aws_cost_forecast_monthly`
  - `aws_cost_usage`
  - `aws_ec2_application_load_balancer`
  - `aws_ec2_autoscaling_group`
  - `aws_ec2_capacity_reservation`
  - `aws_ec2_classic_load_balancer`
  - `aws_ec2_gateway_load_balancer`
  - `aws_ec2_key_pair`
  - `aws_s3_access_point`
  - `aws_s3_account_settings`
  - `aws_vpc`
  - `aws_vpc_customer_gateway`
  - `aws_vpc_dhcp_options`
  - `aws_vpc_eip`
  - `aws_vpc_endpoint`
  - `aws_vpc_flow_log`
  - `aws_vpc_nat_gateway`
  - `aws_vpc_network_acl`
  - `aws_vpc_peering_connection`
  - `aws_vpc_route_table`
  - `aws_vpc_security_group`
  - `aws_vpc_subnet`
  - `aws_vpc_vpn_connection`
  - `aws_vpc_vpn_gateway`
- Updated the query headers in the `aws_api_gatewayv2_api` table documentation.

_Bug fixes_

- Queries will no longer fail if the `regions` config arg is set to `["*"]` when AWS releases a new region that is not included in the plugin's region list. ([#1267](https://github.com/turbot/steampipe-plugin-aws/pull/1267))
- Queries will no longer fail if the `regions` config arg includes a wildcarded item, e.g., `["test-*"]`, that matches on no valid regions. ([#1276](https://github.com/turbot/steampipe-plugin-aws/pull/1276))

## v0.74.2 [2022-08-26]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.4](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v414-2022-08-26) which fixes the query timeout issues during dashboard execution and compliance checks. ([#1264](https://github.com/turbot/steampipe-plugin-aws/pull/1264))

## v0.74.1 [2022-08-25]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v412-2022-08-25) which fixes the stalling of dashboard queries and compliance checks. ([#1259](https://github.com/turbot/steampipe-plugin-aws/pull/1259))

_Bug fixes_

- Fixed the plugin credential caching issue wherein the sessions which had an error were also cached. ([#1255](https://github.com/turbot/steampipe-plugin-aws/pull/1255))

## v0.74.0 [2022-08-24]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v411-2022-08-24) which includes several caching and memory management improvements. ([#1252](https://github.com/turbot/steampipe-plugin-aws/pull/1252))
- Recompiled plugin with Go version `1.19`. ([#1250](https://github.com/turbot/steampipe-plugin-aws/pull/1250))

_What's new?_

- New tables added
  - [aws_dynamodb_table_export](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dynamodb_table_export) ([#1218](https://github.com/turbot/steampipe-plugin-aws/pull/1218))
  - [aws_eks_node_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_node_group) ([#1236](https://github.com/turbot/steampipe-plugin-aws/pull/1236))
  - [aws_emr_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_instance) ([#1225](https://github.com/turbot/steampipe-plugin-aws/pull/1225))
  - [aws_emr_instance_fleet](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_instance_fleet) ([#1226](https://github.com/turbot/steampipe-plugin-aws/pull/1226))

_Enhancements_

- Added column `cluster_arn` to `aws_ecs_container_instance` table. ([#1239](https://github.com/turbot/steampipe-plugin-aws/pull/1239))
- Added column `streaming_destination` to `aws_dynamodb_table` table. ([#1227](https://github.com/turbot/steampipe-plugin-aws/pull/1227))
- Added column `vault_notification_config` to `aws_glacier_vault` table. ([#1231](https://github.com/turbot/steampipe-plugin-aws/pull/1231))
- Added column `file_system_configs` to `aws_lambda_function` table. ([#1224](https://github.com/turbot/steampipe-plugin-aws/pull/1224))

_Bug fixes_

- List queries for the `aws_emr_instance_group` table no longer fail if there are any instance groups in clusters that use instance fleets. ([#1228](https://github.com/turbot/steampipe-plugin-aws/pull/1228))

## v0.73.0 [2022-08-16]

_Enhancements_

- Added column `subnet_id` to `aws_ec2_network_interface` table. ([#1216](https://github.com/turbot/steampipe-plugin-aws/pull/1216))

_Bug fixes_

- Fixed the `aws_eventbridge_rule` table to also list rules for non-default EventBridge buses. ([#1214](https://github.com/turbot/steampipe-plugin-aws/pull/1214))
- Fixed the `aws_rds_db_cluster` table to also list MySQL and PostgreSQL engine type clusters. ([#1213](https://github.com/turbot/steampipe-plugin-aws/pull/1213))

## v0.72.0 [2022-08-15]

_What's new?_

- New tables added
  - [aws_ses_domain_identity](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ses_domain_identity) ([#1206](https://github.com/turbot/steampipe-plugin-aws/pull/1206)) (Thanks to [@janritter](https://github.com/janritter) for the contribution!)

_Enhancements_

- Re-enabled `name` and `type` optional list key columns in `aws_route53_record` table. ([#1190](https://github.com/turbot/steampipe-plugin-aws/pull/1190))
- Updated the following tables to use [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2): ([#1186](https://github.com/turbot/steampipe-plugin-aws/pull/1186))
  - `aws_acm_certificate`
  - `aws_api_gateway_api_authorizer`
  - `aws_api_gateway_api_key`
  - `aws_api_gateway_rest_api`
  - `aws_api_gateway_stage`
  - `aws_api_gateway_usage_plan`
  - `aws_api_gatewayv2_api`
  - `aws_api_gatewayv2_domain_name`
  - `aws_api_gatewayv2_integration`
  - `aws_api_gatewayv2_stage`
  - `aws_dynamodb_backup`
  - `aws_iam_access_advisor`
  - `aws_iam_access_key`
  - `aws_iam_account_password_policy`
  - `aws_iam_account_summary`
  - `aws_iam_credential_report`
  - `aws_iam_group`
  - `aws_iam_policy`
  - `aws_iam_policy_attachment`
  - `aws_iam_policy_simulator`
  - `aws_iam_role`
  - `aws_iam_saml_provider`
  - `aws_iam_server_certificate`
  - `aws_iam_user`
  - `aws_iam_virtual_mfa_device`
  - `aws_s3_bucket`
  - `aws_sns_topic`

_Bug fixes_

- `aws_backup_vault` table now returns no rows instead of an error when querying a vault that does not exist. ([#1163](https://github.com/turbot/steampipe-plugin-aws/pull/1163))
- `aws_neptune_db_cluster` table now only lists Neptune DB clusters. ([#1204](https://github.com/turbot/steampipe-plugin-aws/pull/1204))
- `aws_rds_db_cluster` table now only lists RDS Aurora DB clusters. ([#1204](https://github.com/turbot/steampipe-plugin-aws/pull/1204))

## v0.71.0 [2022-07-20]

_What's new?_

- New tables added
  - [aws_networkfirewall_firewall_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_networkfirewall_firewall_policy) ([#1171](https://github.com/turbot/steampipe-plugin-aws/pull/1171))

_Enhancements_

- Added the following new columns to `aws_sns_topic` table: ([#1176](https://github.com/turbot/steampipe-plugin-aws/pull/1176))
  - application_failure_feedback_role_arn
  - application_success_feedback_role_arn
  - application_success_feedback_sample_rate
  - firehose_failure_feedback_role_arn
  - firehose_success_feedback_role_arn
  - firehose_success_feedback_sample_rate
  - http_failure_feedback_role_arn
  - http_success_feedback_role_arn
  - http_success_feedback_sample_rate
  - lambda_failure_feedback_role_arn
  - lambda_success_feedback_role_arn
  - lambda_success_feedback_sample_rate
  - sqs_failure_feedback_role_arn
  - sqs_success_feedback_role_arn
  - sqs_success_feedback_sample_rate
- Added support for `us-iso` and `us-isob` regions. ([#1168](https://github.com/turbot/steampipe-plugin-aws/pull/1168))

_Bug fixes_

- Fixed the typo in column name to use `health_check_target` instead of `heath_check_target` in `aws_ec2_classic_load_balancer` table. ([#1179](https://github.com/turbot/steampipe-plugin-aws/pull/1179))
- Fixed the `settings` column in the `aws_ecs_cluster` table to correctly return data instead of `null`. ([#1175](https://github.com/turbot/steampipe-plugin-aws/pull/1175))

## v0.70.0 [2022-07-14]

_What's new?_

- New tables added
  - [aws_waf_rule_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_waf_rule_group) ([#1160](https://github.com/turbot/steampipe-plugin-aws/pull/1160))
  - [aws_waf_web_acl](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_waf_web_acl) ([#1151](https://github.com/turbot/steampipe-plugin-aws/pull/1151))

_Enhancements_

- Added column `associated_resources` to `aws_wafv2_web_acl` table. ([#1158](https://github.com/turbot/steampipe-plugin-aws/pull/1158))

## v0.69.0 [2022-07-12]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11). ([#1150](https://github.com/turbot/steampipe-plugin-aws/pull/1150))
- Recompiled plugin with [aws-sdk-go v1.44.49](https://github.com/aws/aws-sdk-go/blob/main/CHANGELOG.md#release-v14449-2022-07-06). ([#1142](https://github.com/turbot/steampipe-plugin-aws/pull/1142))
- Added timestamps to example queries in `aws_cloudtrail_trail_event`, `aws_cloudwatch_log_event` and `aws_vpc_flow_log_event` table documents. ([#1136](https://github.com/turbot/steampipe-plugin-aws/pull/1136))
- Added column `url_config` to `aws_lambda_alias` and `aws_lambda_function` tables. ([#1146](https://github.com/turbot/steampipe-plugin-aws/pull/1146))

_Bug fixes_

- Fixed inconsistent table names in the `aws_ebs_volume_metric_write_ops`, `aws_ebs_volume_metric_write_ops_hourly` and `aws_vpc_flow_log` tables. ([#1149](https://github.com/turbot/steampipe-plugin-aws/pull/1149))

## v0.68.0 [2022-07-06]

_What's new?_

- New tables added
  - [aws_cloudfront_response_headers_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudfront_response_headers_policy) ([#1128](https://github.com/turbot/steampipe-plugin-aws/pull/1128))
  - [aws_iam_saml_provider](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iam_saml_provider) ([#1125](https://github.com/turbot/steampipe-plugin-aws/pull/1125))
  - [aws_pricing_service_attribute](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_pricing_service_attribute) ([#1137](https://github.com/turbot/steampipe-plugin-aws/pull/1137)) ([#1141](https://github.com/turbot/steampipe-plugin-aws/pull/1141))

_Enhancements_

- Added column `certificate` to `aws_rds_db_instance` table. ([#1126](https://github.com/turbot/steampipe-plugin-aws/pull/1126))

_Bug fixes_

- Fixed the `aws_backup_framework` table to return an empty row for the unsupported `ap-northeast-3` region instead of returning an error. ([#1131](https://github.com/turbot/steampipe-plugin-aws/pull/1131))

## v0.67.0 [2022-07-01]

_What's new?_

- New tables added
  - [aws_amplify_app](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_amplify_app) ([#1112](https://github.com/turbot/steampipe-plugin-aws/pull/1112))
  - [aws_cloudfront_function](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudfront_function) ([#1120](https://github.com/turbot/steampipe-plugin-aws/pull/1120))
  - [aws_glue_connection](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_connection) ([#1102](https://github.com/turbot/steampipe-plugin-aws/pull/1102))
  - [aws_glue_data_catalog_encryption_settings](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_data_catalog_encryption_settings) ([#1114](https://github.com/turbot/steampipe-plugin-aws/pull/1114))
  - [aws_glue_job](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_job) ([#1118](https://github.com/turbot/steampipe-plugin-aws/pull/1118))
  - [aws_glue_security_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_security_configuration) ([#1106](https://github.com/turbot/steampipe-plugin-aws/pull/1106))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v331--2022-06-30). ([#1129](https://github.com/turbot/steampipe-plugin-aws/pull/1129))
- Added information about STS and IAM API limitations with aws-vault temporary credentials in the `docs/index.md` file.
- Added column `vpcs` to `aws_route53_zone` table. ([#1085](https://github.com/turbot/steampipe-plugin-aws/pull/1085))
- Added column `vpc_endpoint_service_permissions` to `aws_vpc_endpoint_service` table. ([#1121](https://github.com/turbot/steampipe-plugin-aws/pull/1121))

_Bug fixes_

- Fixed the `No such host` issue in audit manager tables. ([#1122](https://github.com/turbot/steampipe-plugin-aws/pull/1122))
- Fixed the `MaxResults` parameter issue in list API for `aws_eks_identity_provider_config` table. ([#1119](https://github.com/turbot/steampipe-plugin-aws/pull/1119))
- Fixed the `Unsupported region` issue in `aws_media_store_container` table. ([#1117](https://github.com/turbot/steampipe-plugin-aws/pull/1117))
- Fixed the `BdRequestException` issue in the `aws_guardduty_member` table. ([#1116](https://github.com/turbot/steampipe-plugin-aws/pull/1116))

## v0.66.0 [2022-06-24]

_What's new?_

- New tables added
  - [aws_backup_framework](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_framework`) ([#1099](https://github.com/turbot/steampipe-plugin-aws/pull/1099))
  - [aws_elasticache_reserved_cache_node](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_reserved_cache_node) ([#1092](https://github.com/turbot/steampipe-plugin-aws/pull/1092))
- Added `s3_force_path_style` config argument to allow S3 path-style addressing. ([#1082](https://github.com/turbot/steampipe-plugin-aws/pull/1082)) (Thanks to [@srgg](https://github.com/srgg) for the contribution!)

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v330--2022-6-22). ([#1108](https://github.com/turbot/steampipe-plugin-aws/pull/1108))

## v0.65.0 [2022-06-16]

_What's new?_

- New tables added
  - [aws_rds_reserved_db_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_reserved_db_instance) ([#1087](https://github.com/turbot/steampipe-plugin-aws/pull/1087))

_Enhancements_

- Added column `pending_maintenance_actions` to `aws_rds_db_cluster` and `aws_rds_db_instance` tables. ([#1083](https://github.com/turbot/steampipe-plugin-aws/pull/1083))
- Updated the `.gitignore` file to include all VS Code user settings. ([#1078](https://github.com/turbot/steampipe-plugin-aws/pull/1078))

_Bug fixes_

- Fixed the `snapshot_create_time` column in `aws_redshift_snapshot` table to be of `timestamp` data type instead of `string`. ([#1071](https://github.com/turbot/steampipe-plugin-aws/pull/1071))

## v0.64.0 [2022-06-09]

_What's new?_

- New tables added
  - [aws_elasticache_redis_metric_engine_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_engine_cpu_utilization_daily) ([#1063](https://github.com/turbot/steampipe-plugin-aws/pull/1063))
  - [aws_glue_dev_endpoint](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_dev_endpoint) ([#1057](https://github.com/turbot/steampipe-plugin-aws/pull/1057))
  - [aws_ssm_inventory](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_inventory) ([#1054](https://github.com/turbot/steampipe-plugin-aws/pull/1054))

_Enhancements_

- Updated `aws_route53_record` table to temporarily disable `name` and `type` list key quals in order to fix duplicate rows issue. ([#972](https://github.com/turbot/steampipe-plugin-aws/pull/972))

_Bug fixes_

- Fixed `aws_elasticsearch_domain`, `aws_opensearch_domain`, and `aws_s3_bucket` tables to not panic when ignoring errors. ([#1064](https://github.com/turbot/steampipe-plugin-aws/pull/1064))

## v0.63.0 [2022-06-03]

_What's new?_

- Added `endpoint_url` config arg to provide users the ability to set a custom endpoint URL when making requests to AWS services. For more information, please see [AWS plugin configuration](https://hub.steampipe.io/plugins/turbot/aws#configuration). ([#1053](https://github.com/turbot/steampipe-plugin-aws/pull/1053)) (Thanks to [@srgg](https://github.com/srgg) for the contribution!)

## v0.62.0 [2022-06-02]

_What's new?_

- New tables added
  - [aws_route53_traffic_policy_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_traffic_policy_instance) ([#1047](https://github.com/turbot/steampipe-plugin-aws/pull/1047))

_Enhancements_

- Added column `administrator_account` to `aws_securityhub_hub` table. ([#1046](https://github.com/turbot/steampipe-plugin-aws/pull/1046))

_Bug fixes_

- Fixed the `is_logging` column of `aws_cloudtrail_trail` table to return `true` instead of `null` for shadow trails when the source trail has logging enabled. ([#986](https://github.com/turbot/steampipe-plugin-aws/pull/986))

## v0.61.0 [2022-05-30]

_What's new?_

- New tables added
  - [aws_route53_health_check](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_health_check) ([#1045](https://github.com/turbot/steampipe-plugin-aws/pull/1045))

_Bug fixes_

- Fixed the `inline_policies` column in `aws_iam_role`, `aws_iam_group` and `aws_iam_user` tables to correctly return results instead of an error. ([#1048](https://github.com/turbot/steampipe-plugin-aws/pull/1048))

## v0.60.0 [2022-05-25]

_What's new?_

- Added `ignore_error_codes` config arg to provide users the ability to set a list of additional AWS error codes to ignore while running queries. For instance, to ignore some common access denied errors, which is helpful when running with limited permissions, set the argument `ignore_error_codes = ["AccessDenied", "AccessDeniedException"]`. For more information, please see [AWS plugin configuration](https://hub.steampipe.io/plugins/turbot/aws#configuration) ([#992](https://github.com/turbot/steampipe-plugin-aws/pull/992))
- New tables added
  - [aws_config_aggregate_authorization](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_config_aggregate_authorization) ([#1025](https://github.com/turbot/steampipe-plugin-aws/pull/1025))
  - [aws_dlm_lifecycle_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dlm_lifecycle_policy) ([#1016](https://github.com/turbot/steampipe-plugin-aws/pull/1016))
  - [aws_guardduty_filter](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_filter) ([#1029](https://github.com/turbot/steampipe-plugin-aws/pull/1029))
  - [aws_guardduty_member](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_member) ([#1028](https://github.com/turbot/steampipe-plugin-aws/pull/1028))
  - [aws_guardduty_publishing_destination](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_publishing_destination) ([#1030](https://github.com/turbot/steampipe-plugin-aws/pull/1030))
  - [aws_inspector_assessment_run](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector_assessment_run) ([#1036](https://github.com/turbot/steampipe-plugin-aws/pull/1036))
  - [aws_inspector_exclusion](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector_exclusion) ([#1038](https://github.com/turbot/steampipe-plugin-aws/pull/1038))
  - [aws_inspector_finding](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector_finding) ([#1040](https://github.com/turbot/steampipe-plugin-aws/pull/1040))
  - [aws_ram_resource_association](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ram_resource_association) ([#1009](https://github.com/turbot/steampipe-plugin-aws/pull/1009))
  - [aws_ram_principal_association](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ram_principal_association) ([#1009](https://github.com/turbot/steampipe-plugin-aws/pull/1009))
  - [aws_securityhub_action_target](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_action_target) ([#1012](https://github.com/turbot/steampipe-plugin-aws/pull/1012))
  - [aws_securityhub_finding](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_finding) ([#1017](https://github.com/turbot/steampipe-plugin-aws/pull/1017))
  - [aws_securityhub_finding_aggregator](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_finding_aggregator) ([#1031](https://github.com/turbot/steampipe-plugin-aws/pull/1031))
  - [aws_securityhub_insight](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_insight) ([#1011](https://github.com/turbot/steampipe-plugin-aws/pull/1011))
  - [aws_securityhub_member](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_member) ([#1022](https://github.com/turbot/steampipe-plugin-aws/pull/1022))
  - [aws_securityhub_standards_control](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_standards_control) ([#1010](https://github.com/turbot/steampipe-plugin-aws/pull/1010))

_Enhancements_

- Added column `shared_directories` to `aws_directory_service_directory` table. ([#1024](https://github.com/turbot/steampipe-plugin-aws/pull/1024))
- Added column `vpc_id` to `aws_ec2_network_interface` table. ([#990](https://github.com/turbot/steampipe-plugin-aws/pull/990))
- Added column `master_account` to `aws_guardduty_detector` table. ([#1023](https://github.com/turbot/steampipe-plugin-aws/pull/1023))
- Added column `architectures` to `aws_lambda_function` table. ([#991](https://github.com/turbot/steampipe-plugin-aws/pull/991))
- Updated all tables to use `IgnoreConfig` instead of `ShouldIgnoreError` in `GetConfig` function. ([#992](https://github.com/turbot/steampipe-plugin-aws/pull/992))

_Bug fixes_

- Fixed the handling for unsupported regions in `aws_inspector_assessment_target` and `aws_inspector_assessment_template` tables. ([#1039](https://github.com/turbot/steampipe-plugin-aws/pull/1039)

## v0.59.0 [2022-05-11]

_What's new?_

- New tables added
  - [aws_glue_catalog_table](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_catalog_table) ([#963](https://github.com/turbot/steampipe-plugin-aws/pull/963))
  - [aws_opensearch_domain](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_opensearch_domain) ([#984](https://github.com/turbot/steampipe-plugin-aws/pull/984))
  - [aws_pinpoint_app](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_pinpoint_app) ([#968](https://github.com/turbot/steampipe-plugin-aws/pull/968))
  - [aws_route53_traffic_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_traffic_policy) ([#983](https://github.com/turbot/steampipe-plugin-aws/pull/983))
  - [aws_sagemaker_app](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sagemaker_app) ([#977](https://github.com/turbot/steampipe-plugin-aws/pull/977))
  - [aws_ses_email_identity](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ses_email_identity) ([#978](https://github.com/turbot/steampipe-plugin-aws/pull/978))

_Enhancements_

- Improved the example descriptions in `aws_iam_credential_report` table document.

_Bug fixes_

- Fixed `aws_cloudtrail_trail_event`, `aws_cloudwatch_log_event`, and `aws_vpc_flow_log_event` tables not returning correct results for consecutive queries when using the `filter` list key column. ([#981](https://github.com/turbot/steampipe-plugin-aws/pull/981))

## v0.58.0 [2022-05-05]

_What's new?_

- New tables added
  - [aws_neptune_db_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_neptune_db_cluster) ([#966](https://github.com/turbot/steampipe-plugin-aws/pull/966))
  - [aws_sagemaker_domain](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sagemaker_domain) ([#974](https://github.com/turbot/steampipe-plugin-aws/pull/974))

_Enhancements_

- Added the `environment_variables` column to `aws_lambda_function` and `aws_lambda_version` tables. ([#973](https://github.com/turbot/steampipe-plugin-aws/pull/973))
- Updated the `aws_organizations_account` table's `id` column description and document for account ID clarifications. ([#975](https://github.com/turbot/steampipe-plugin-aws/pull/975))
- Removed the use of chalk package in `aws_iam_credential_report` table for dashboard compatibility.

_Bug fixes_

- Updated the column name from `date-created` to `date_created` in the `aws_elastic_beanstalk_environment` table ([#965](https://github.com/turbot/steampipe-plugin-aws/pull/965))

## v0.57.0 [2022-04-27]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#953](https://github.com/turbot/steampipe-plugin-aws/pull/953))
- Added support for native Linux ARM and Mac M1 builds. ([#958](https://github.com/turbot/steampipe-plugin-aws/pull/958))
- Added column `package_type` to `aws_lambda_function` table. ([#956](https://github.com/turbot/steampipe-plugin-aws/pull/956))

## v0.56.0 [2022-04-13]

_What's new?_

- New tables added
  - [aws_cost_by_record_type_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_record_type_daily) ([#950](https://github.com/turbot/steampipe-plugin-aws/pull/950))
  - [aws_cost_by_record_type_monthly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_record_type_monthly) ([#950](https://github.com/turbot/steampipe-plugin-aws/pull/950))
  - [aws_wafregional_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafregional_rule) ([#931](https://github.com/turbot/steampipe-plugin-aws/pull/931))
- Added optional config arguments `max_error_retry_attempts` and `min_error_retry_delay` to allow customization of the error retry timings. For more information please see [AWS plugin configuration](https://hub.steampipe.io/plugins/turbot/aws#configuration). ([#914](https://github.com/turbot/steampipe-plugin-aws/pull/914))

_Enhancements_

- Added column `event_notification_configuration` to `aws_s3_bucket` table. ([#946](https://github.com/turbot/steampipe-plugin-aws/pull/946))
- Added column `login_profile` to `aws_iam_user` table. ([#947](https://github.com/turbot/steampipe-plugin-aws/pull/947))

## v0.55.0 [2022-04-06]

_Enhancements_

- Added `image_scanning_findings` column to `aws_ecr_repository` table ([#937](https://github.com/turbot/steampipe-plugin-aws/pull/937))

## v0.54.0 [2022-04-01]

- New tables added
  - [aws_networkfirewall_rule_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_networkfirewall_rule_group) ([#944](https://github.com/turbot/steampipe-plugin-aws/pull/944))

## v0.53.0 [2022-03-30]

_Enhancements_

- Added `table_class` column to `aws_dynamodb_table` table ([#936](https://github.com/turbot/steampipe-plugin-aws/pull/936))
- Added additional optional key quals ('!=') to `aws_cost_by_service_daily`, `aws_cost_by_service_monthly`, `aws_cost_by_service_usage_type_daily` and `aws_cost_by_service_usage_type_monthly` tables and context cancellation to `aws_cost_forecast_daily` and `aws_cost_forecast_monthly` tables ([#917](https://github.com/turbot/steampipe-plugin-aws/pull/917))

_Bug fixes_

- Fixed `aws_s3_bucket` queries failing for buckets created in the `EU` (eu-west-1) region through the CLI or API ([#927](https://github.com/turbot/steampipe-plugin-aws/pull/927))

## v0.52.0 [2022-03-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#933](https://github.com/turbot/steampipe-plugin-aws/pull/933))

## v0.51.0 [2022-03-17]

_Enhancements_

- Added column `standards_status_reason_code` to `aws_securityhub_standards_subscription` table ([#930](https://github.com/turbot/steampipe-plugin-aws/pull/930))

## v0.50.1 [2022-03-10]

_Bug fixes_

- Fixed the `aws_ebs_snapshot` table to correctly handle `InvalidParameterValue` error ([#919](https://github.com/turbot/steampipe-plugin-aws/pull/919))

## v0.50.0 [2022-03-04]

_Enhancements_

- Added `sqs_managed_sse_enabled` column to `aws_sqs_queue` table ([#922](https://github.com/turbot/steampipe-plugin-aws/pull/922))
- Added additional optional key quals to `aws_cost_by_service_daily`, `aws_cost_by_service_monthly`, `aws_cost_by_service_usage_type_daily` and `aws_cost_by_service_usage_type_monthly` tables ([#912](https://github.com/turbot/steampipe-plugin-aws/pull/912))

_Bug fixes_

- Fixed the `title` column of `aws_vpc_security_group_rule` table to correctly evaluate if a security group rule is either ingress or egress ([#924](https://github.com/turbot/steampipe-plugin-aws/pull/924))

## v0.49.0 [2022-02-17]

_What's new?_

- New tables added
  - [aws_servicequotas_service_quota_change_request](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicequotas_service_quota_change_request) ([#889](https://github.com/turbot/steampipe-plugin-aws/pull/889))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.0.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v203--2022-02-14) ([#901](https://github.com/turbot/steampipe-plugin-aws/pull/901))

_Bug fixes_

- Fixed pagination issues in ` aws_ecs_service` table ([#908](https://github.com/turbot/steampipe-plugin-aws/pull/908))
- Fixed the `aws_iam_access_advisor` table to handle the errors when steampipe is running on multi-account connections by using an [aggregator connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators) in the configuration file ([#905](https://github.com/turbot/steampipe-plugin-aws/pull/905))

## v0.48.0 [2022-02-14]

_What's new?_

- New tables added
  - [aws_cloudwatch_metric](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_metric) ([#880](https://github.com/turbot/steampipe-plugin-aws/pull/880))

_Enhancements_

- Added context cancellation handling to the following tables ([#896](https://github.com/turbot/steampipe-plugin-aws/pull/896))
  - aws_auditmanager_control
  - aws_auditmanager_framework
  - aws_backup_recovery_point
  - aws_backup_vault
  - aws_cloudfront_cache_policy
  - aws_cloudtrail_trail
  - aws_cloudtrail_trail_event
  - aws_cloudwatch_log_event
  - aws_cloudwatch_log_resource_policy
  - aws_ec2_reserved_instance
  - aws_guardduty_finding
  - aws_iam_action
  - aws_kinesis_video_stream
  - aws_lambda_alias
  - aws_lambda_function
  - aws_serverlessapplicationrepository_application
  - aws_ssm_patch_baseline
  - aws_vpc_security_group_rule

- Updated default max records parameter value and lower limit for the following tables ([#896](https://github.com/turbot/steampipe-plugin-aws/pull/896))
  - aws_api_gateway_api_authorizer
  - aws_api_gatewayv2_stage
  - aws_config_conformance_pack
  - aws_directory_service_directory
  - aws_ecs_container_instance
  - aws_ecs_service

_Bug fixes_

- Fixed the `aws_codecommit_repository` table to correctly list out all the repositories ([#894](https://github.com/turbot/steampipe-plugin-aws/pull/894))

## v0.47.0 [2022-02-09]

_What's new?_

- New tables added
  - [aws_glue_crawler](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_crawler) ([#882](https://github.com/turbot/steampipe-plugin-aws/pull/882))
  - [aws_servicequotas_default_service_quota](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicequotas_default_service_quota) ([#887](https://github.com/turbot/steampipe-plugin-aws/pull/887))
  - [aws_servicequotas_service_quota](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_servicequotas_service_quota) ([#884](https://github.com/turbot/steampipe-plugin-aws/pull/884))

## v0.46.0 [2022-02-02]

_Enhancements_

- Added additional optional key quals, filter support, and context cancellation handling to `Redshift`, `Route 53`, `S3`, `SageMaker`, `Secrets Manager`, `Security Hub`, `Serverless Application Repository`, `Step Functions`, `SNS`, `SSM`, `SSO`, `VPC`, `WAF` and `Well-Architected` tables ([#873](https://github.com/turbot/steampipe-plugin-aws/pull/873))

_Bug fixes_

- Fixed the `aws_dax_cluster` table to skip unsupported regions ([#869](https://github.com/turbot/steampipe-plugin-aws/pull/869))
- Fixed the `aws_wellarchitected_workload` table to skip unsupported regions ([#859](https://github.com/turbot/steampipe-plugin-aws/pull/859))
- Fixed the `aws_vpc_security_group_rule` table to set the `pair_group_name` column to `nil` for cross-account referenced security group rules instead of returning an error ([#875](https://github.com/turbot/steampipe-plugin-aws/pull/875))
- Updated the column type of `created_date` and `last_modified` columns to `TIMESTAMP` in all Lambda tables ([#871](https://github.com/turbot/steampipe-plugin-aws/pull/871))

## v0.45.0 [2022-01-28]

_Enhancements_

- Added additional optional key quals, filter support, and context cancellation handling to `FSx`, `Glacier`, `GuardDuty`, `IAM`, `Identity Store`, `Inspector`, `Kinesis`, `KMS`, `Lambda`, `ElastiCache`, `Macie` and `RDS` tables ([#856](https://github.com/turbot/steampipe-plugin-aws/pull/856))
- Added the following columns to the `aws_vpc_security_group_rule` table ([#860](https://github.com/turbot/steampipe-plugin-aws/pull/860))
  - cidr_ipv4
  - description
  - group_owner_id
  - is_egress
  - referenced_group_id
  - referenced_peering_status
  - referenced_user_id
  - referenced_vpc_id
  - referenced_vpc_peering_connection_id
  - security_group_rule_id
- Added `assignment_status` column to `aws_iam_virtual_mfa_device` table ([#856](https://github.com/turbot/steampipe-plugin-aws/pull/856))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#865](https://github.com/turbot/steampipe-plugin-aws/pull/865))

_Bug fixes_

- Fixed the `aws_workspaces_workspace` table to skip the unsupported regions ([#862](https://github.com/turbot/steampipe-plugin-aws/pull/862))

_Deprecated_

- The following columns of `aws_vpc_security_group_rule` table have been deprecated to stay consistent with the API response data. These columns will be removed in the next major version. We recommend updating any scripts or workflows that use these deprecated columns to use the equivalent new columns in the table instead.
  - cidr_ip (replaced by cidr_ipv4)
  - group_name
  - owner_id (replaced by group_owner_id)
  - pair_group_id (replaced by referenced_group_id)
  - pair_group_name
  - pair_peering_status (replaced by referenced_peering_status)
  - pair_user_id (replaced by referenced_user_id)
  - pair_vpc_id (replaced by referenced_vpc_id)
  - pair_vpc_peering_connection_id (replaced by referenced_vpc_peering_connection_id)
  - vpc_id

## v0.44.0 [2022-01-12]

_Enhancements_

- Recompiled plugin with [aws-sdk-go-v1.42.25](https://github.com/aws/aws-sdk-go/blob/main/CHANGELOG.md#release-v14225-2021-12-21) ([#851](https://github.com/turbot/steampipe-plugin-aws/pull/851))
- Added additional optional key quals, filter support, and context cancellation handling to `ACM`, `API Gateway`, `EBS`, `EC2`, `ECR`, `ECS`, `EFS`, `EKS`, `Elastic Beanstalk`, `ElastiCache`, `Elasticsearch`, `EMR`, `EventBridge` and `RDS` tables ([#850](https://github.com/turbot/steampipe-plugin-aws/pull/850))

## v0.43.0 [2021-12-21]

_What's new?_

- New tables added
  - [aws_iam_policy_attachment](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iam_policy_attachment) ([#824](https://github.com/turbot/steampipe-plugin-aws/pull/824))

_Enhancements_

- Updated default max records parameter value and lower limit for `Access Analyzer`, `ACM`, `API Gateway`, `Application Auto Scaling`, `Audit manager`, `Backup`, `Cloud Control`, `CloudFormation`, `CloudFront`, `CloudWatch`, `CodePipeline`, `Config`, `DAX`, `DMS` and `DynamoDB` tables ([#829](https://github.com/turbot/steampipe-plugin-aws/pull/829))

_Bug fixes_

- Fixed the `aws_workspaces_workspace` table to return an empty row for unsupported regions instead of throwing an error ([#835](https://github.com/turbot/steampipe-plugin-aws/pull/835))
- Querying the `aws_ebs_snapshot` table will now correctly return snapshot(s) details instead of an empty row ([#842](https://github.com/turbot/steampipe-plugin-aws/pull/842))
- The `image_owner_alias` column of `aws_ec2_ami_shared` table will now correctly display the AWS account alias (for example, amazon, self) or the AWS account ID of the AMI owner ([#841](https://github.com/turbot/steampipe-plugin-aws/pull/841))
- The `image_owner_alias` column of `aws_ec2_ami` table is now set to `self` by default ([#841](https://github.com/turbot/steampipe-plugin-aws/pull/841))

## v0.42.2 [2021-12-14]

_Bug fixes_

- Fixed default max records parameter value and lower limit for `aws_cloudwatch_alarm` table

## v0.42.1 [2021-12-14]

_Bug fixes_

- Queries no longer fail when using a wildcard in the `region` config argument due to the release of `ap-southeast-3` region

## v0.42.0 [2021-12-08]

_What's new?_

- New tables added
  - [aws_ec2_managed_prefix_list](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_managed_prefix_list) ([#813](https://github.com/turbot/steampipe-plugin-aws/pull/813))
  - [aws_vpc_peering_connection](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_peering_connection) ([#814](https://github.com/turbot/steampipe-plugin-aws/pull/814))

_Enhancements_

- Added column `prefix_list_id` to `aws_vpc_security_group_rule` table ([#801](https://github.com/turbot/steampipe-plugin-aws/pull/801))
- Added column `compliance_by_config_rule` to `aws_config_rule` table ([#817](https://github.com/turbot/steampipe-plugin-aws/pull/817))
- Added column `project_visibility` to `aws_codebuild_project` table ([821](https://github.com/turbot/steampipe-plugin-aws/pull/821))
- Added additional optional key quals, filter support, and context cancellation handling to `Access Analyzer`, `ACM`, `API Gateway`, `Application Auto Scaling`, `Audit manager`, `Backup`, `Cloud Control`, `CloudFormation`, `CloudFront`, `CloudWatch`, `CodeBuild`, `CodeCommit`, `CodePipeline`, `Config`, `DAX`, `Directory Service`, `DMS`, `DynamoDB` and `EBS` tables ([754](https://github.com/turbot/steampipe-plugin-aws/pull/754))
- Added an example query for listing SQL server instances with SSL disabled in the `aws_rds_db_instance` table document ([#806](https://github.com/turbot/steampipe-plugin-aws/pull/806))
- `README.md` and `docs/index.md` files now have updated Slack channel links

_Bug fixes_

- Fixed the `string field contains invalid UTF-8` error in the `aws_ec2_instance` table ([#812](https://github.com/turbot/steampipe-plugin-aws/pull/812))

## v0.41.0 [2021-11-23]

_What's new?_

- New tables added
  - [aws_elasticache_redis_metric_cache_hits_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_cache_hits_hourly) ([#753](https://github.com/turbot/steampipe-plugin-aws/pull/753))
  - [aws_elasticache_redis_metric_curr_connections_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_curr_connections_hourly) ([#753](https://github.com/turbot/steampipe-plugin-aws/pull/753))
  - [aws_elasticache_redis_metric_engine_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_engine_cpu_utilization_hourly) ([#753](https://github.com/turbot/steampipe-plugin-aws/pull/753))
  - [aws_elasticache_redis_metric_get_type_cmds_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_get_type_cmds_hourly) ([#753](https://github.com/turbot/steampipe-plugin-aws/pull/753))
  - [aws_elasticache_redis_metric_list_based_cmds_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_list_based_cmds_hourly) ([#753](https://github.com/turbot/steampipe-plugin-aws/pull/753))
  - [aws_elasticache_redis_metric_new_connections_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_redis_metric_new_connections_hourly) ([#753](https://github.com/turbot/steampipe-plugin-aws/pull/753))
  - [aws_serverlessapplicationrepository_application](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_serverlessapplicationrepository_application) ([#751](https://github.com/turbot/steampipe-plugin-aws/pull/751))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#741](https://github.com/turbot/steampipe-plugin-aws/pull/795))
- Added filter example queries in `aws_cloudwatch_log_event` table ([#748](https://github.com/turbot/steampipe-plugin-aws/pull/748))
- Added few more example queries in `aws_iam_role` table ([#685](https://github.com/turbot/steampipe-plugin-aws/pull/685))

_Bug fixes_

- `aws_ec2_application_load_balancer` table will no longer return `ValidationError` in get call ([#792](https://github.com/turbot/steampipe-plugin-aws/pull/792))
- `aws_dax_cluster` table will no longer return an error when we try to query for unsupported regions ([#787](https://github.com/turbot/steampipe-plugin-aws/pull/787))
- `aws_lambda_alias` table will now need `name`, `function_name` and `region` to perform get call ([#781](https://github.com/turbot/steampipe-plugin-aws/pull/781))

## v0.40.0 [2021-11-17]

_What's new?_

- New tables added
  - [aws_cloudwatch_log_resource_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_resource_policy) ([#747](https://github.com/turbot/steampipe-plugin-aws/pull/747))
  - [aws_media_store_container](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_media_store_container) ([#749](https://github.com/turbot/steampipe-plugin-aws/pull/749))

_Enhancements_

- Updated: Add `policy_std` column to the `aws_ecrpublic_repository` table ([#778](https://github.com/turbot/steampipe-plugin-aws/pull/778))
- Updated: Add `policy_std` column to the `aws_ecr_repository` table ([#780](https://github.com/turbot/steampipe-plugin-aws/pull/780))
- Updated: Add columns `policy` and `policy_std` to the `aws_lambda_alias` table ([#774](https://github.com/turbot/steampipe-plugin-aws/pull/774))
- Updated: Add columns `policy` and `policy_std` to the `aws_lambda_version` table ([#776](https://github.com/turbot/steampipe-plugin-aws/pull/776))
- Updated: Add columns `policy` and `policy_std` to the `aws_secretsmanager_secret` table ([#745](https://github.com/turbot/steampipe-plugin-aws/pull/745))

_Bug fixes_

- Fixed: `aws_kinesis_firehose_delivery_stream` table now includes better error handling ([#769](https://github.com/turbot/steampipe-plugin-aws/pull/769))
- Fixed: Remove duplicate data from the `aws_backup_plan` table ([#767](https://github.com/turbot/steampipe-plugin-aws/pull/767))
- Fixed: `aws_ecrpublic_repository` table will now return an empty row instead of an error when we try to query for any region other than `us-east-1` ([#770](https://github.com/turbot/steampipe-plugin-aws/pull/770))

## v0.39.1 [2021-11-15]

_Bug fixes_

- Fixed: Queries will no longer panic when encountering an error due to invalid references in the `ShouldRetry` function ([#763](https://github.com/turbot/steampipe-plugin-aws/pull/763))

## v0.39.0 [2021-11-12]

_What's new?_

- New tables added
  - [aws_eventbridge_bus](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eventbridge_bus) ([#737](https://github.com/turbot/steampipe-plugin-aws/pull/737))
  - [aws_lambda_layer](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lambda_layer) ([#740](https://github.com/turbot/steampipe-plugin-aws/pull/740))
  - [aws_lambda_layer_version](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lambda_layer_version) ([#743](https://github.com/turbot/steampipe-plugin-aws/pull/743))

_Enhancements_

- Updated: Add `policy_std` column to `aws_backup_vault` table ([#746](https://github.com/turbot/steampipe-plugin-aws/pull/746))
- Updated: Increase the golangci-lint workflow timeout to 10 mins ([#750](https://github.com/turbot/steampipe-plugin-aws/pull/750))

_Bug fixes_

- Fixed: Queries will no longer hang if no credentials are provided or an invalid profile is specified ([#713](https://github.com/turbot/steampipe-plugin-aws/pull/713))

## v0.38.2 [2021-11-08]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.7.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v173--2021-11-08) ([#741](https://github.com/turbot/steampipe-plugin-aws/pull/741))

## v0.38.1 [2021-11-05]

_Bug fixes_

- Updated data type of the column `platform_version` from `int` to `string` in `aws_ssm_managed_instance` table ([#732](https://github.com/turbot/steampipe-plugin-aws/pull/732))

## v0.38.0 [2021-11-03]

_What's new?_

- New tables added
  - [aws_ecs_task](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_task) ([#715](https://github.com/turbot/steampipe-plugin-aws/pull/715))
  - [aws_sfn_state_machine](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sfn_state_machine) ([#714](https://github.com/turbot/steampipe-plugin-aws/pull/714))
  - [aws_sfn_state_machine_execution](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sfn_state_machine_execution) ([#724](https://github.com/turbot/steampipe-plugin-aws/pull/724))
  - [aws_sfn_state_machine_execution_history](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sfn_state_machine_execution_history) ([#728](https://github.com/turbot/steampipe-plugin-aws/pull/728))

_Enhancements_

- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.7.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v171--2021-11-01) ([#729](https://github.com/turbot/steampipe-plugin-aws/pull/729))

## v0.37.0 [2021-10-27]

_What's new?_

- New tables added
  - [aws_backup_protected_resource](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_protected_resource) ([#704](https://github.com/turbot/steampipe-plugin-aws/pull/704))
  - [aws_backup_recovery_point](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_recovery_point) ([#705](https://github.com/turbot/steampipe-plugin-aws/pull/705))
  - [aws_fsx_file_system](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_fsx_file_system) ([#693](https://github.com/turbot/steampipe-plugin-aws/pull/693))
  - [aws_identitystore_user](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_identitystore_user) ([#675](https://github.com/turbot/steampipe-plugin-aws/pull/675))

_Enhancements_

- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.7.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v170--2021-10-18) ([#695](https://github.com/turbot/steampipe-plugin-aws/pull/695))

_Bug fixes_

- Queries for global tables, e.g., `aws_iam_user`, will no longer return an error if no regions are specified for a connection ([#690](https://github.com/turbot/steampipe-plugin-aws/pull/690))
- Fixed the `ecs_service` table to correctly return the tags instead of returning `null` ([#710](https://github.com/turbot/steampipe-plugin-aws/pull/710))

## v0.36.0 [2021-10-12]

_What's new?_

- New tables added
  - [aws_ssoadmin_managed_policy_attachment](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssoadmin_managed_policy_attachment) ([#664](https://github.com/turbot/steampipe-plugin-aws/pull/664))
  - [aws_workspaces_workspace](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_workspaces_workspace) ([#681](https://github.com/turbot/steampipe-plugin-aws/pull/681))

## v0.35.1 [2021-10-08]

_Bug fixes_

- Fixed: Increase number of retries from 3->8 for Cloud Control service sessions to better handle throttling
- Fixed: Examples for `aws_cloudcontrol_resource` table are now correct

## v0.35.0 [2021-10-08]

_What's new?_

- New tables added
  - [aws_cloudcontrol_resource](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudcontrol_resource) ([#680](https://github.com/turbot/steampipe-plugin-aws/pull/680))
  - [aws_ssoadmin_permission_set](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssoadmin_permission_set) ([#659](https://github.com/turbot/steampipe-plugin-aws/pull/659))

_Enhancements_

- Updated: Parliament IAM permissions to the latest ([#676](https://github.com/turbot/steampipe-plugin-aws/pull/676))
- Updated: Add additional optional key quals, filter support, and context cancellation handling to `aws_ec2_instance`, `aws_iam_policy`, `aws_rds_db_cluster_snapshot` tables ([#638](https://github.com/turbot/steampipe-plugin-aws/pull/638))
- Recompiled plugin with [steampipe-plugin-sdk v1.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v162--2021-10-08)

## v0.34.0 [2021-09-30]

_What's new?_

- New tables added
  - [aws_identitystore_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_identitystore_group) ([#663](https://github.com/turbot/steampipe-plugin-aws/pull/663))

_Bug fixes_

- Add pagination to list and list tags functions in several tables ([#660](https://github.com/turbot/steampipe-plugin-aws/pull/660))

## v0.33.0 [2021-09-22]

_What's new?_

- New tables added
  - [aws_organizations_account](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_organizations_account) ([#650](https://github.com/turbot/steampipe-plugin-aws/pull/650))
  - [aws_ssoadmin_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssoadmin_instance) ([#658](https://github.com/turbot/steampipe-plugin-aws/pull/658))

_Bug fixes_

- When the macie service is not enabled in a particular region, `aws_macie2_classification_job` table will now return `nil` instead of `error` ([#661](https://github.com/turbot/steampipe-plugin-aws/pull/661))

## v0.32.1 [2021-09-13]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v151--2021-09-13) ([#653](https://github.com/turbot/steampipe-plugin-aws/pull/653))

## v0.32.0 [2021-09-09]

_What's new?_

- New tables added
  - [aws_ec2_transit_gateway_route](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_transit_gateway_route) ([#637](https://github.com/turbot/steampipe-plugin-aws/pull/637))

_Enhancements_

- Added customized exponential back-off retry logic to optimize retry mechanism ([#635](https://github.com/turbot/steampipe-plugin-aws/pull/635))

_Bug fixes_

- Fixed: Implemented pagination in `aws_config_rule` and `aws_config_conformance_pack` table ([#646](https://github.com/turbot/steampipe-plugin-aws/pull/646))
- Fixed: Improved documentations ([#639](https://github.com/turbot/steampipe-plugin-aws/pull/639))

## v0.31.0 [2021-08-25]

_What's new?_

- New tables added
  - [aws_ec2_capacity_reservation](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_capacity_reservation) ([#619](https://github.com/turbot/steampipe-plugin-aws/pull/619))
  - [aws_ec2_reserved_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_reserved_instance) ([#612](https://github.com/turbot/steampipe-plugin-aws/pull/612))
  - [aws_lambda_function_metric_duration_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lambda_function_metric_duration_daily) ([#624](https://github.com/turbot/steampipe-plugin-aws/pull/624))
  - [aws_lambda_function_metric_errors_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lambda_function_metric_errors_daily) ([#625](https://github.com/turbot/steampipe-plugin-aws/pull/625))
  - [aws_lambda_function_metric_invocations_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_lambda_function_metric_invocations_daily) ([#627](https://github.com/turbot/steampipe-plugin-aws/pull/627))
  - [aws_tagging_resource](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_tagging_resource) ([#628](https://github.com/turbot/steampipe-plugin-aws/pull/628))

## v0.30.0 [2021-08-20]

_What's new?_

- New tables added
  - [aws_directory_service_directory](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_directory_service_directory) ([#572](https://github.com/turbot/steampipe-plugin-aws/pull/572))
  - [aws_ec2_application_load_balancer_metric_request_count_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_application_load_balancer_metric_request_count_daily) ([#605](https://github.com/turbot/steampipe-plugin-aws/pull/605))
  - [aws_ec2_network_load_balancer_metric_net_flow_count](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_network_load_balancer_metric_net_flow_count) ([#527](https://github.com/turbot/steampipe-plugin-aws/pull/527))
  - [aws_ec2_application_load_balancer_metric_request_count](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_application_load_balancer_metric_request_count) ([#527](https://github.com/turbot/steampipe-plugin-aws/pull/527))
  - [aws_ec2_network_load_balancer_metric_net_flow_count_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_network_load_balancer_metric_net_flow_count_daily) ([#604](https://github.com/turbot/steampipe-plugin-aws/pull/604))
  - [aws_rds_db_event_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_event_subscription) ([#609](https://github.com/turbot/steampipe-plugin-aws/pull/609))
  - [aws_redshift_cluster_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshift_cluster_metric_cpu_utilization_daily) ([#606](https://github.com/turbot/steampipe-plugin-aws/pull/606))
  - [aws_securityhub_standards_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_standards_subscription) ([#286](https://github.com/turbot/steampipe-plugin-aws/pull/286))

_Enhancements_

- Updated: Global services like IAM, S3, Route 53, etc. will now connect to `us-gov-west-1` and `cn-northwest-1` when creating service connections in GovCloud and China respectively ([#613](https://github.com/turbot/steampipe-plugin-aws/pull/613))
- Updated: Add column `scheduled_actions` to `aws_redshift_cluster` table ([#523](https://github.com/turbot/steampipe-plugin-aws/pull/523))
- Updated: Add column `log_publishing_options` to `aws_elasticsearch_domain` table ([#593](https://github.com/turbot/steampipe-plugin-aws/pull/593))
- Updated: Add column `instance_lifecycle` to `aws_ec2_instance` table ([#617](https://github.com/turbot/steampipe-plugin-aws/pull/617))

_Bug fixes_

- Fixed: `aws_ec2_ssl_policy` table will no longer generate duplicate values with multi-region setup ([#594](https://github.com/turbot/steampipe-plugin-aws/pull/594))
- Fixed: If no regions are set in the config file, the region will now correctly be determined from the `AWS_DEFAULT_REGION` or `AWS_REGION` environment variables if set ([#598](https://github.com/turbot/steampipe-plugin-aws/pull/598))

## v0.29.0 [2021-08-06]

_What's new?_

- New tables added
  - [aws_eks_identity_provider_config](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_identity_provider_config) ([#571](https://github.com/turbot/steampipe-plugin-aws/pull/571))

_Bug fixes_

- Fixed: `arn` column data now contain the correct regions in regional resource tables ([#590](https://github.com/turbot/steampipe-plugin-aws/pull/590))
- Fixed: Querying columns `dnssec_key_signing_keys` and `dnssec_status` in `aws_route53_zone` table for private hosted zones no longer causes errors ([#589](https://github.com/turbot/steampipe-plugin-aws/pull/589))

## v0.28.0 [2021-08-05]

_What's new?_

- New tables added
  - [aws_emr_cluster_metric_is_idle](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_cluster_metric_is_idle) ([#570](https://github.com/turbot/steampipe-plugin-aws/pull/570))

_Bug fixes_

- Fixed: `aws_cloudtrail_trail` table is now smarter when hydrating data for shadow trails (global and organization) ([#578](https://github.com/turbot/steampipe-plugin-aws/pull/578))
- Fixed: Route tables with IPv6 routes no longer cause queries to fail in the `aws_vpc_route` table ([#581](https://github.com/turbot/steampipe-plugin-aws/pull/581))

## v0.27.0 [2021-07-31]

_What's new?_

- New tables added
  - [aws_cloudtrail_trail_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_trail_event) ([#564](https://github.com/turbot/steampipe-plugin-aws/pull/564))
  - [aws_cloudwatch_log_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_event) ([#564](https://github.com/turbot/steampipe-plugin-aws/pull/564))
  - [aws_ecs_cluster_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_cluster_metric_cpu_utilization) ([#563](https://github.com/turbot/steampipe-plugin-aws/pull/563))
  - [aws_ecs_cluster_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_cluster_metric_cpu_utilization_daily) ([#563](https://github.com/turbot/steampipe-plugin-aws/pull/563))
  - [aws_ecs_cluster_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_cluster_metric_cpu_utilization_hourly) ([#563](https://github.com/turbot/steampipe-plugin-aws/pull/563))
  - [aws_emr_instance_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_instance_group) ([#562](https://github.com/turbot/steampipe-plugin-aws/pull/562))
  - [aws_vpc_flow_log_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log_event) ([#564](https://github.com/turbot/steampipe-plugin-aws/pull/564))

_Bug fixes_

- Fixed: `aws_ec2_instance` table should not panic when hydrating `state_transition_time` column if there is no state transition reason ([#574](https://github.com/turbot/steampipe-plugin-aws/pull/574))

## v0.26.0 [2021-07-22]

_What's new?_

- New tables added
  - [aws_codecommit_repository](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codecommit_repository) ([#515](https://github.com/turbot/steampipe-plugin-aws/pull/515))
  - [aws_ecs_service](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_service) ([#555](https://github.com/turbot/steampipe-plugin-aws/pull/555))

_Enhancements_

- Updated: Add column `arn` in `aws_vpc_nat_gateway` table ([#540](https://github.com/turbot/steampipe-plugin-aws/pull/540))
- Updated: Add multi-account connection information and examples to index doc ([#565](https://github.com/turbot/steampipe-plugin-aws/pull/565))
- Updated: Improve error message when connection config regions are not valid ([#558](https://github.com/turbot/steampipe-plugin-aws/pull/558))
- Updated: Cleanup region selection in connection creation code for table modules ([#566](https://github.com/turbot/steampipe-plugin-aws/pull/566))
- Recompiled plugin with [steampipe-plugin-sdk v1.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v141--2021-07-20)

_Bug fixes_

- Fixed: Connection creation is now retried when receiving reset by peer errors ([#557](https://github.com/turbot/steampipe-plugin-aws/pull/557))
- Fixed: Fix plugin sometimes incorrectly selecting the wrong region from connection config ([#561](https://github.com/turbot/steampipe-plugin-aws/pull/561))
- Fixed: Hydration now works for `created_at`, `name`, `title`, `updated_at`, and `version` columns in `aws_codepipeline_pipeline` table ([#537](https://github.com/turbot/steampipe-plugin-aws/pull/537))
- Fixed: Several column descriptions in `aws_ecs_task_definition` table ([#541](https://github.com/turbot/steampipe-plugin-aws/pull/541))

## v0.25.0 [2021-07-08]

_What's new?_

- New tables added
  - [aws_auditmanager_evidence](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_auditmanager_evidence) ([#450](https://github.com/turbot/steampipe-plugin-aws/pull/450))
  - [aws_config_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_config_rule) ([#417](https://github.com/turbot/steampipe-plugin-aws/pull/417))
  - [aws_dynamodb_metric_account_provisioned_read_capacity_util](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dynamodb_metric_account_provisioned_read_capacity_util) ([#518](https://github.com/turbot/steampipe-plugin-aws/pull/518))
  - [aws_dynamodb_metric_account_provisioned_write_capacity_util](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dynamodb_metric_account_provisioned_write_capacity_util) ([#518](https://github.com/turbot/steampipe-plugin-aws/pull/518))

_Enhancements_

- Updated: Add wildcard support when defining regions in plugin connection configuration ([#530](https://github.com/turbot/steampipe-plugin-aws/pull/530))
- Updated: Improve docs/index.md with expanded credential options and examples ([#535](https://github.com/turbot/steampipe-plugin-aws/pull/535))

_Bug fixes_

- Fixed: Fix various failing integration tests ([#534](https://github.com/turbot/steampipe-plugin-aws/pull/534))
- Fixed: Removed invalid key column definitions in various tables

## v0.24.0 [2021-07-01]

_What's new?_

- New tables added
  - [aws_codebuild_source_credential](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codebuild_source_credential) ([#511](https://github.com/turbot/steampipe-plugin-aws/pull/511))
  - [aws_macie2_classification_job](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_macie2_classification_job) ([#479](https://github.com/turbot/steampipe-plugin-aws/pull/479))

_Enhancements_

- Updated: Rename column `file_system_arn` to `arn` in `aws_efs_file_system` table ([#494](https://github.com/turbot/steampipe-plugin-aws/pull/494))
- Updated: Rename column `table_arn` to `arn` in `aws_dynamodb_table` table ([#495](https://github.com/turbot/steampipe-plugin-aws/pull/495))
- Updated: Improve error message in `aws_iam_credential_report` table when no credential report exists ([#510](https://github.com/turbot/steampipe-plugin-aws/pull/510))
- Updated: Remove use of deprecated function `ItemFromKey` from `aws_redshift_cluster` table ([#514](https://github.com/turbot/steampipe-plugin-aws/pull/514))

## v0.23.0 [2021-06-24]

_What's new?_

- New tables added
  - [aws_auditmanager_evidence_folder](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_auditmanager_evidence_folder) ([#435](https://github.com/turbot/steampipe-plugin-aws/pull/435))
  - [aws_backup_selection](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_selection) ([#501](https://github.com/turbot/steampipe-plugin-aws/pull/501))
  - [aws_codepipeline_pipeline](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codepipeline_pipeline) ([#498](https://github.com/turbot/steampipe-plugin-aws/pull/498))
  - [aws_route53_domain](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_domain) ([#360](https://github.com/turbot/steampipe-plugin-aws/pull/360))

_Bug fixes_

- Fixed: Typo in description for common cloudwatch_metric `timestamp` column ([#505](https://github.com/turbot/steampipe-plugin-aws/pull/505))

## v0.22.0 [2021-06-17]

_What's new?_

- New tables added
  - [aws_ec2_ssl_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_ssl_policy) ([#362](https://github.com/turbot/steampipe-plugin-aws/pull/362))
  - [aws_eks_addon_version](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_addon_version) ([#482](https://github.com/turbot/steampipe-plugin-aws/pull/482))
  - [aws_guardduty_finding](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_finding) ([#488](https://github.com/turbot/steampipe-plugin-aws/pull/488))
  - [aws_ssm_managed_instance_compliance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_managed_instance_compliance) ([#484](https://github.com/turbot/steampipe-plugin-aws/pull/484))
  - [aws_vpc_vpn_connection](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_vpn_connection) ([#486](https://github.com/turbot/steampipe-plugin-aws/pull/486))

_Enhancements_

- Updated: Add column `arn` to `aws_api_gateway_stage` table ([#447](https://github.com/turbot/steampipe-plugin-aws/pull/447))
- Updated: Add column `arn` to `aws_ec2_classic_load_balancer` table ([#475](https://github.com/turbot/steampipe-plugin-aws/pull/475))
- Updated: Add column `event_subscriptions` to `aws_inspector_assessment_template` table ([#467](https://github.com/turbot/steampipe-plugin-aws/pull/467))
- Updated: Add column `logging_configuration` to `aws_wafv2_web_acl` table ([#470](https://github.com/turbot/steampipe-plugin-aws/pull/470))
- Updated: Add columns `dnssec_key_signing_keys` and `dnssec_status` to `aws_route53_zone` table ([#439](https://github.com/turbot/steampipe-plugin-aws/pull/439))

_Bug fixes_

- Fixed: Cache key in `ElasticsearchService` function and update various cache keys to be more consistent ([#500](https://github.com/turbot/steampipe-plugin-aws/pull/500))
- Fixed: Tags hydrate call should not fail in `aws_sagemaker_notebook_instance` table ([#372](https://github.com/turbot/steampipe-plugin-aws/pull/372))

## v0.21.0 [2021-06-10]

_What's new?_

- New tables added
  - [aws_ec2_ami_shared](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_ami_shared) ([#456](https://github.com/turbot/steampipe-plugin-aws/pull/456))
  - [aws_ec2_regional_settings](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_regional_settings) ([#403](https://github.com/turbot/steampipe-plugin-aws/pull/403))
  - [aws_eks_addon](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_addon) ([#478](https://github.com/turbot/steampipe-plugin-aws/pull/478))
  - [aws_sagemaker_endpoint_configuration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sagemaker_endpoint_configuration) ([#477](https://github.com/turbot/steampipe-plugin-aws/pull/477))

_Enhancements_

- Updated: Shadow trails are now included in `aws_cloudtrail_trail` table query results ([#441](https://github.com/turbot/steampipe-plugin-aws/pull/441))
- Updated: Add columns `replication_group_id`, `snapshot_retention_limit`, and `snapshot_window` to `aws_elasticache_cluster` table ([#458](https://github.com/turbot/steampipe-plugin-aws/pull/458))
- Updated: Add columns `dead_letter_config_target_arn` and `reserved_concurrent_executions` to `aws_lambda_function` table ([#474](https://github.com/turbot/steampipe-plugin-aws/pull/474))
- Updated: Rename column `alarm_arn` to `arn` in `aws_cloudwatch_alarm` table ([#489](https://github.com/turbot/steampipe-plugin-aws/pull/489))
- Recompiled plugin with [steampipe-plugin-sdk v0.2.10](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v0210-2021-06-09)

_Bug fixes_

- Fixed: GetCommonColumns function should only get STS caller identity once per account instead of per region ([#490](https://github.com/turbot/steampipe-plugin-aws/pull/490))

## v0.20.0 [2021-06-03]

_What's new?_

- New tables added
  - [aws_auditmanager_assessment](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_auditmanager_assessment) ([#430](https://github.com/turbot/steampipe-plugin-aws/pull/430))
  - [aws_auditmanager_control](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_auditmanager_control) ([#425](https://github.com/turbot/steampipe-plugin-aws/pull/425))
  - [aws_auditmanager_framework](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_auditmanager_framework) ([#415](https://github.com/turbot/steampipe-plugin-aws/pull/415))
  - [aws_cloudfront_cache_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudfront_cache_policy) ([#431](https://github.com/turbot/steampipe-plugin-aws/pull/431))
  - [aws_cloudfront_origin_access_identity](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudfront_origin_access_identity) ([#427](https://github.com/turbot/steampipe-plugin-aws/pull/427))
  - [aws_cloudfront_origin_request_policy](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudfront_origin_request_policy) ([#432](https://github.com/turbot/steampipe-plugin-aws/pull/432))

_Enhancements_

- Updated: Add column `arn` to `aws_redshift_cluster` table ([#462](https://github.com/turbot/steampipe-plugin-aws/pull/462))
- Updated: Add column `arn` to `aws_vpc_network_acl` table ([#457](https://github.com/turbot/steampipe-plugin-aws/pull/457))
- Updated: Add column `object_lock_configuration` to `aws_s3_bucket` table ([#464](https://github.com/turbot/steampipe-plugin-aws/pull/464))
- Updated: Add column `state_transition_time` to `aws_ec2_instance` table ([#344](https://github.com/turbot/steampipe-plugin-aws/pull/344))
- Updated: Bump urllib3 in /scripts/generate_parliament_iam_permissions ([#471](https://github.com/turbot/steampipe-plugin-aws/pull/471))
- Updated: Getting tags for clusters in 'creating' state should not error in `aws_elasticache_cluster` table ([#454](https://github.com/turbot/steampipe-plugin-aws/pull/454))
- Updated: Rename column `replication_instance_arn` to `arn` in `aws_dms_replication_instance` table ([#455](https://github.com/turbot/steampipe-plugin-aws/pull/455))

_Bug fixes_

- Fixed: Rename `table_ aws_elasticsearch_domain.go` to `table_aws_elasticsearch_domain.go`

## v0.19.0 [2021-05-27]

_What's new?_

- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)
- New tables added
  - [aws_api_gatewayv2_integration](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_api_gatewayv2_integration) ([#346](https://github.com/turbot/steampipe-plugin-aws/pull/346))
  - [aws_cloudfront_distribution](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudfront_distribution) ([#388](https://github.com/turbot/steampipe-plugin-aws/pull/388))
  - [aws_cost_by_account_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_account_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_by_account_monthly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_account_monthly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_by_service_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_service_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_by_service_monthly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_service_monthly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_by_service_usage_type_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_service_usage_type_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_by_service_usage_type_monthly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_by_service_usage_type_monthly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_forecast_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_forecast_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_forecast_monthly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_forecast_monthly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_cost_usage](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cost_usage) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ebs_volume_metric_read_ops](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ebs_volume_metric_read_ops) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ebs_volume_metric_read_ops_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ebs_volume_metric_read_ops_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ebs_volume_metric_read_ops_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ebs_volume_metric_read_ops_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ebs_volume_metric_write_ops](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ebs_volume_metric_write_ops) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ebs_volume_metric_write_ops_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ebs_volume_metric_write_ops_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ebs_volume_metric_write_ops_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ebs_volume_metric_write_ops_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ec2_instance_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_instance_metric_cpu_utilization) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ec2_instance_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_instance_metric_cpu_utilization_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_ec2_instance_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_instance_metric_cpu_utilization_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_efs_mount_target](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_efs_mount_target) ([#426](https://github.com/turbot/steampipe-plugin-aws/pull/426))
  - [aws_kinesisanalyticsv2_application](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_kinesisanalyticsv2_application) ([#358](https://github.com/turbot/steampipe-plugin-aws/pull/358))
  - [aws_rds_db_instance_metric_connections](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_connections) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_connections_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_connections_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_connections_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_connections_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_cpu_utilization) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_cpu_utilization_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_cpu_utilization_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_read_iops](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_read_iops) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_read_iops_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_read_iops_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_read_iops_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_read_iops_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_write_iops](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_write_iops) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_write_iops_daily](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_write_iops_daily) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_rds_db_instance_metric_write_iops_hourly](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_rds_db_instance_metric_write_iops_hourly) ([#437](https://github.com/turbot/steampipe-plugin-aws/pull/437))
  - [aws_sagemaker_training_job](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sagemaker_training_job) ([#384](https://github.com/turbot/steampipe-plugin-aws/pull/384))
  - [aws_ssm_managed_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_managed_instance) ([#436](https://github.com/turbot/steampipe-plugin-aws/pull/436))
  - [aws_waf_rate_based_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_waf_rate_based_rule) ([#289](https://github.com/turbot/steampipe-plugin-aws/pull/289))
  - [aws_wafv2_rule_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafv2_rule_group) ([#281](https://github.com/turbot/steampipe-plugin-aws/pull/281))

_Enhancements_

- Updated: Base64 data in the `user_data` column is now decoded in the `aws_ec2_instance` and `aws_ec2_launch_configuration` tables ([#363](https://github.com/turbot/steampipe-plugin-aws/pull/363))
- Updated: Add `arn` column to `aws_account` table ([#418](https://github.com/turbot/steampipe-plugin-aws/pull/418))
- Updated: Add `arn` column to `aws_guardduty_detector` table ([#408](https://github.com/turbot/steampipe-plugin-aws/pull/408))
- Updated: Add `arn` column to `aws_ssm_association` table ([#404](https://github.com/turbot/steampipe-plugin-aws/pull/404))

## v0.18.0 [2021-05-20]

_What's new?_

- New tables added
  - [aws_glue_catalog_database](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glue_catalog_database) ([#337](https://github.com/turbot/steampipe-plugin-aws/pull/337))

_Enhancements_

- Updated: Add `arn` column to `aws_ebs_snapshot` table ([#405](https://github.com/turbot/steampipe-plugin-aws/pull/405))
- Updated: Add `arn` column to `aws_vpc_eip` table ([#407](https://github.com/turbot/steampipe-plugin-aws/pull/407))
- Updated: Improve availability zone count example in `aws_lambda_function` table doc ([#413](https://github.com/turbot/steampipe-plugin-aws/pull/413))

_Bug fixes_

- Fixed: Getting key rotation status for external keys should not error in `aws_kms_key` table ([#398](https://github.com/turbot/steampipe-plugin-aws/pull/398))

## v0.17.0 [2021-05-13]

_What's new?_

- New tables added
  - [aws_dms_replication_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dms_replication_instance) ([#357](https://github.com/turbot/steampipe-plugin-aws/pull/357))
  - [aws_ecs_container_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_container_instance) ([#340](https://github.com/turbot/steampipe-plugin-aws/pull/340))
  - [aws_sagemaker_model](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sagemaker_model) ([#371](https://github.com/turbot/steampipe-plugin-aws/pull/371))
  - [aws_waf_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_waf_rule) ([#287](https://github.com/turbot/steampipe-plugin-aws/pull/287))
  - [aws_wafv2_ip_set](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafv2_ip_set) ([#255](https://github.com/turbot/steampipe-plugin-aws/pull/255))
  - [aws_wafv2_regex_pattern_set](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafv2_regex_pattern_set) ([#276](https://github.com/turbot/steampipe-plugin-aws/pull/276))

_Enhancements_

- Updated: README.md and docs/index.md now contain links to our Slack community ([#411](https://github.com/turbot/steampipe-plugin-aws/pull/411))
- Updated: Add `logging_status` column to `aws_redshift_cluster` table ([#350](https://github.com/turbot/steampipe-plugin-aws/pull/350))
- Updated: Add missing columns available in the hydrate data to `aws_ssm_association` table ([#356](https://github.com/turbot/steampipe-plugin-aws/pull/356))
- Updated: Bump lodash from 4.17.20 to 4.17.21 in /aws-test ([#389](https://github.com/turbot/steampipe-plugin-aws/pull/389))

_Bug fixes_

- Fixed: Querying the aws_iam_account_password_policy table should not error if no password policy exists ([#382](https://github.com/turbot/steampipe-plugin-aws/pull/382))

## v0.16.0 [2021-05-06]

_What's new?_

- New tables added
  - [aws_accessanalyzer_analyzer](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_accessanalyzer_analyzer) ([#341](https://github.com/turbot/steampipe-plugin-aws/pull/341))
  - [aws_appautoscaling_target](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_appautoscaling_target) ([#353](https://github.com/turbot/steampipe-plugin-aws/pull/353))
  - [aws_wafv2_web_acl](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wafv2_web_acl) ([#245](https://github.com/turbot/steampipe-plugin-aws/pull/245))

_Enhancements_

- Updated: Add `arn` column to `aws_config_configuration_recorder` table ([#380](https://github.com/turbot/steampipe-plugin-aws/pull/380))
- Updated: Add `arn` column to `aws_ebs_volume` table ([#368](https://github.com/turbot/steampipe-plugin-aws/pull/368))
- Updated: Add `arn` column to `aws_ec2_instance` table ([#367](https://github.com/turbot/steampipe-plugin-aws/pull/367))
- Updated: Add `arn` column to `aws_vpc_security_group` table ([#377](https://github.com/turbot/steampipe-plugin-aws/pull/377))
- Updated: Add `arn` column to `aws_vpc` table ([#378](https://github.com/turbot/steampipe-plugin-aws/pull/378))
- Updated: Add `automatic_backups` column to `aws_efs_file_system` table ([#351](https://github.com/turbot/steampipe-plugin-aws/pull/351))

_Bug fixes_

- Fixed: Handling of pending subscriptions in `aws_sns_topic_subscription` table ([#349](https://github.com/turbot/steampipe-plugin-aws/pull/349))

## v0.15.0 [2021-04-29]

_What's new?_

- New tables added
  - [aws_dax_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_dax_cluster) ([#328](https://github.com/turbot/steampipe-plugin-aws/pull/328))
  - [aws_ecrpublic_repository](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecrpublic_repository) ([#336](https://github.com/turbot/steampipe-plugin-aws/pull/336))
  - [aws_s3_access_point](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_s3_access_point) ([#318](https://github.com/turbot/steampipe-plugin-aws/pull/318))

_Enhancements_

- Updated: Parliament IAM permissions for Parliament v1.4.0 ([#216](https://github.com/turbot/steampipe-plugin-aws/issues/216))

_Bug fixes_

- Fixed: The `aws_guardduty_threat_intel_set` table should not throw an rpc error while trying to list threat intel sets ([#343](https://github.com/turbot/steampipe-plugin-aws/pull/343))

## v0.14.0 [2021-04-22]

_What's new?_

- New tables added
  - [aws_codebuild_project](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_codebuild_project) ([#319](https://github.com/turbot/steampipe-plugin-aws/pull/319))
  - [aws_elasticsearch_domain](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticsearch_domain) ([#327](https://github.com/turbot/steampipe-plugin-aws/pull/327))
  - [aws_sagemaker_notebook_instance](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_sagemaker_notebook_instance) ([#324](https://github.com/turbot/steampipe-plugin-aws/pull/324))
  - [aws_secretsmanager_secret](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_secretsmanager_secret) ([#330](https://github.com/turbot/steampipe-plugin-aws/pull/330))

_Bug fixes_

- Fixed: Replace hardcoded ARN references in `aws_ec2_instance_type`, `aws_iam_policy`, and `aws_s3_bucket` tables ([#331](https://github.com/turbot/steampipe-plugin-aws/pull/331))

## v0.13.0 [2021-04-15]

_What's new?_

- New tables added
  - [aws_elasticache_subnet_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_subnet_group) ([#247](https://github.com/turbot/steampipe-plugin-aws/pull/247))
  - [aws_guardduty_detector](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_detector) ([#251](https://github.com/turbot/steampipe-plugin-aws/pull/251))
  - [aws_guardduty_ipset](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_ipset) ([#259](https://github.com/turbot/steampipe-plugin-aws/pull/259))
  - [aws_guardduty_threat_intel_set](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_guardduty_threat_intel_set) ([#271](https://github.com/turbot/steampipe-plugin-aws/pull/271))
  - [aws_inspector_assessment_template](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector_assessment_template) ([#248](https://github.com/turbot/steampipe-plugin-aws/pull/248))
  - [aws_redshift_snapshot](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshift_snapshot) ([#238](https://github.com/turbot/steampipe-plugin-aws/pull/238))

_Enhancements_

- Updated: Add `arn` column to `aws_s3_bucket` table ([#313](https://github.com/turbot/steampipe-plugin-aws/pull/313))

_Bug fixes_

- Fixed: Query example in `aws_iam_server_certificate` table docs ([#309](https://github.com/turbot/steampipe-plugin-aws/pull/309))

## v0.12.0 [2021-04-08]

_What's new?_

- New tables added
  - [aws_inspector_assessment_target](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_inspector_assessment_target) ([#234](https://github.com/turbot/steampipe-plugin-aws/pull/234))

_Enhancements_

- Updated: Add `metadata_options` column to `aws_ec2_instance` table ([#306](https://github.com/turbot/steampipe-plugin-aws/pull/306))

## v0.11.0 [2021-04-08]

_What's new?_

- New tables added
  - [aws_backup_plan](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_plan) ([#164](https://github.com/turbot/steampipe-plugin-aws/pull/164))
  - [aws_backup_vault](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_backup_vault) ([#163](https://github.com/turbot/steampipe-plugin-aws/pull/163))
  - [aws_elastic_beanstalk_application](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elastic_beanstalk_application) ([#187](https://github.com/turbot/steampipe-plugin-aws/pull/187))
  - [aws_iam_server_certificate](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iam_server_certificate) ([#293](https://github.com/turbot/steampipe-plugin-aws/pull/293))
  - [aws_redshift_event_subscription](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshift_event_subscription) ([#230](https://github.com/turbot/steampipe-plugin-aws/pull/230))
  - [aws_securityhub_product](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_hub) ([#207](https://github.com/turbot/steampipe-plugin-aws/pull/207))
  - [aws_wellarchitected_workload](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_wellarchitected_workload) ([#223](https://github.com/turbot/steampipe-plugin-aws/pull/223))

_Enhancements_

- Updated: Add `certificate_transparency_logging_preference`, `imported_at`, `renewal_eligibility`, and `type` columns to `aws_acm_certificate` table ([#299](https://github.com/turbot/steampipe-plugin-aws/pull/299))
- Recompiled plugin with [steampipe-plugin-sdk v0.2.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v027-2021-03-31)

_Bug fixes_

- Fixed: Rename column `instance_profile_arn` to `instance_profile_arns` and update data to be a list of ARNs (strings) in `aws_iam_role` table ([#291](https://github.com/turbot/steampipe-plugin-aws/pull/291))
- Fixed: Release dates in CHANGELOG no longer project versions out in the year 20201 ([#284](https://github.com/turbot/steampipe-plugin-aws/pull/284))

## v0.10.1 [2021-04-02]

_Bug fixes_

- Fixed: `Table definitions & examples` link now points to the correct location ([#282](https://github.com/turbot/steampipe-plugin-aws/pull/282))

## v0.10.0 [2021-04-01]

_What's new?_

- New tables added
  - [aws_cloudwatch_alarm](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_alarm) ([#197](https://github.com/turbot/steampipe-plugin-aws/pull/197))
  - [aws_ecr_repository](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecr_repository) ([#139](https://github.com/turbot/steampipe-plugin-aws/pull/139))
  - [aws_ecs_task_definition](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_task_definition) ([#173](https://github.com/turbot/steampipe-plugin-aws/pull/173))
  - [aws_efs_access_point](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_efs_access_point) ([#174](https://github.com/turbot/steampipe-plugin-aws/pull/174))
  - [aws_elastic_beanstalk_environment](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elastic_beanstalk_environment) ([#178](https://github.com/turbot/steampipe-plugin-aws/pull/178))
  - [aws_elasticache_replication_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_replication_group) ([#246](https://github.com/turbot/steampipe-plugin-aws/pull/246))
  - [aws_glacier_vault](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_glacier_vault) ([#165](https://github.com/turbot/steampipe-plugin-aws/pull/165))
  - [aws_kinesis_consumer](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_kinesis_consumer) ([#222](https://github.com/turbot/steampipe-plugin-aws/pull/222))
  - [aws_redshift_subnet_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshift_subnet_group) ([#220](https://github.com/turbot/steampipe-plugin-aws/pull/220))
  - [aws_securityhub_hub](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_securityhub_hub) ([#166](https://github.com/turbot/steampipe-plugin-aws/pull/166))

_Enhancements_

- Updated: Add `canary_settings` and `method_settings` columns to `aws_api_gateway_stage` table ([#273](https://github.com/turbot/steampipe-plugin-aws/pull/273))
- Updated: Add `query_logging_configs` column to `aws_route53_zone` table ([#264](https://github.com/turbot/steampipe-plugin-aws/pull/264))
- Updated: Example queries for `aws_s3_bucket` table to be more consistent with standards ([#268](https://github.com/turbot/steampipe-plugin-aws/pull/268))

_Bug fixes_

- Fixed: Remove unnecessary engine and region compatibility check when describing instances in the `aws_rds_db_instance` table ([#263](https://github.com/turbot/steampipe-plugin-aws/pull/263))
- Fixed: The `aws_vpc` table should ignore `InvalidVpcID.NotFound` errors ([#270](https://github.com/turbot/steampipe-plugin-aws/pull/270))

## v0.9.0 [2021-03-25]

_What's new?_

- New tables added
  - [aws_config_conformance_pack](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_config_conformance_pack) ([#170](https://github.com/turbot/steampipe-plugin-aws/pull/170))
  - [aws_ecs_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ecs_cluster) ([#140](https://github.com/turbot/steampipe-plugin-aws/pull/140))
  - [aws_efs_file_system](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_efs_file_system) ([#144](https://github.com/turbot/steampipe-plugin-aws/pull/144))
  - [aws_elasticache_parameter_group](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_parameter_group) ([#176](https://github.com/turbot/steampipe-plugin-aws/pull/176))
  - [aws_emr_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_emr_cluster) ([#152](https://github.com/turbot/steampipe-plugin-aws/pull/152))
  - [aws_kinesis_video_stream](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_kinesis_video_stream) ([#182](https://github.com/turbot/steampipe-plugin-aws/pull/182))
  - [aws_route53_resolver_endpoint](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_resolver_endpoint) ([#137](https://github.com/turbot/steampipe-plugin-aws/pull/137))
  - [aws_route53_resolver_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_route53_resolver_rule) ([#148](https://github.com/turbot/steampipe-plugin-aws/pull/148))

_Enhancements_

- Updated: Add `flow_log_status` column to `aws_vpc_flow_log` table ([#233](https://github.com/turbot/steampipe-plugin-aws/pull/233))
- Updated: Add `launch_time` column to `aws_ec2_instance` table ([#227](https://github.com/turbot/steampipe-plugin-aws/pull/227))
- Recompiled plugin with [steampipe-plugin-sdk v0.2.6](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v026-2021-03-18)

## v0.8.0 [2021-03-18]

_What's new?_

- New tables added
  - [aws_cloudtrail_trail](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_trail) ([#34](https://github.com/turbot/steampipe-plugin-aws/pull/34))
  - [aws_eks_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_cluster) ([#131](https://github.com/turbot/steampipe-plugin-aws/pull/131))
  - [aws_elasticache_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_elasticache_cluster) ([#130](https://github.com/turbot/steampipe-plugin-aws/pull/130))
  - [aws_eventbridge_rule](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eventbridge_rule) ([#135](https://github.com/turbot/steampipe-plugin-aws/pull/135))
  - [aws_kinesis_stream](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_kinesis_stream) ([#125](https://github.com/turbot/steampipe-plugin-aws/pull/125))
  - [aws_redshift_cluster](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_redshift_cluster) ([#204](https://github.com/turbot/steampipe-plugin-aws/pull/204))
  - [aws_ssm_association](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_association) ([#114](https://github.com/turbot/steampipe-plugin-aws/pull/114))
  - [aws_ssm_document](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_document) ([#110](https://github.com/turbot/steampipe-plugin-aws/pull/110))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v0.2.4](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v024-2021-03-16)

_Bug fixes_

- Fixed: Various examples for `aws_iam_access_advisor`, `aws_iam_policy_simulator`, and `aws_route53_record` tables ([#186](https://github.com/turbot/steampipe-plugin-aws/pull/186))
- Fixed: Multi-region queries now work properly for the `aws_lambda_version` table ([#192](https://github.com/turbot/steampipe-plugin-aws/pull/192))
- Fixed: `aws_availability_zone` and `aws_ec2_instance_availability` tables now check region opt-in status to avoid `AuthFailure` errors ([#168](https://github.com/turbot/steampipe-plugin-aws/pull/168))
- Fixed: `region` column in `aws_region` table now shows the correct region instead of `global` ([#133](https://github.com/turbot/steampipe-plugin-aws/pull/133))

## v0.7.0 [2021-03-11]

_What's new?_

- New tables added
  - [aws_cloudwatch_log_stream](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudwatch_log_stream)
  - [aws_config_configuration_recorder](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_config_configuration_recorder)
  - [aws_iam_virtual_mfa_device](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_iam_virtual_mfa_device)
  - [aws_ssm_maintenance_window](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_maintenance_window)
  - [aws_ssm_patch_baseline](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ssm_patch_baseline)

_Bug fixes_

- Removed use of deprecated `ItemFromKey` function from all tables

## v0.6.0 [2021-03-05]

_What's new?_

- Plugin now supports authentication through **AWS SSO**.
- New tables added
  - [aws_ec2_gateway_load_balancer](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_ec2_gateway_load_balancer)
  - [aws_vpc_flow_log](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log)

_Enhancements_

- Updated: Added `tags_src` and `tags` columns to `aws_iam_policy` table.

## v0.5.3 [2021-03-02]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve issue:
  - Fix tables failing with error similar to `Error: pq: rpc error: code = Internal desc = get hydrate function getS3Bucket failed with panic interface conversion: interface {} is nil, not *s3.Bucket`([#89](https://github.com/turbot/steampipe-plugin-aws/issues/89)).

## v0.5.2 [2021-02-25]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve SDK issues:
  - Fix error for missing required quals [#40](https://github.com/turbot/steampipe-plugin-sdk/issues/42).
  - Queries fail with error socket: too many open files [#190](https://github.com/turbot/steampipe/issues/190)

## v0.5.1 [2021-02-22]

_Bug fixes_

- Ensure `aws_account` and `aws_region` table work when **regions** argument is specified in connection config ([#70](https://github.com/turbot/steampipe-plugin-aws/pull/70))

## v0.5.0 [2021-02-18]

_What's new?_

- Added support for [connection configuration](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/index.md#connection-configuration). You may specify aws profiles, credentials, and regions for each connection in a configuration file. You can have multiple aws connections, each configured for a different aws account.
- Added multi-region support. A single connection can query multiple AWS regions, via the `regions` connection argument.

_Enhancements_

- Updated: Updated `tag_list` columns to `tags_src` for below RDS service tables.

  - aws_rds_db_cluster
  - aws_rds_db_cluster_parameter_group
  - aws_rds_db_cluster_snapshot
  - aws_rds_db_instance
  - aws_rds_db_option_group
  - aws_rds_db_parameter_group
  - aws_rds_db_snapshot
  - aws_rds_db_subnet_group

- Updated: added `inline_policies_std` column to `aws_iam_group`, `aws_iam_role` and `aws_iam_user` table with canoncialized inline policies.

## v0.4.0 [2021-02-11]

_What's new?_

- New tables added to plugin

  - [aws_iam_access_advisor](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_iam_access_advisor.md) ([#42](https://github.com/turbot/steampipe-plugin-aws/issues/42))
  - [aws_route53_record](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_route53_record.md) ([#43](https://github.com/turbot/steampipe-plugin-aws/issues/43))
  - [aws_route53_zone](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_route53_zone.md) ([#43](https://github.com/turbot/steampipe-plugin-aws/issues/43))

_Enhancements_

- Updated: `aws_iam_credential_report` table to have `password_status` column ([#48](https://github.com/turbot/steampipe-plugin-aws/issues/48))

## v0.3.0 [2021-02-04]

_What's new?_

- New tables added to plugin([#40](https://github.com/turbot/steampipe-plugin-aws/pull/40))

  - [aws_iam_account_password_policy](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_iam_account_password_policy.md)
  - [aws_iam_account_summary](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_iam_account_summary.md)
  - [aws_iam_action](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_iam_action.md)
  - [aws_iam_credential_report](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_iam_credential_report.md)
  - [aws_iam_policy_simulator](https://github.com/turbot/steampipe-plugin-aws/blob/main/docs/tables/aws_iam_policy_simulator.md)

_Enhancements_

- Updated: `aws_ssm_parameter` table to have `value, arn, selector and source_result` fields ([#22](https://github.com/turbot/steampipe-plugin-aws/pull/22))

- Updated: `aws_iam_user` table to have `mfa_enabled and mfa_devices` columns ([#28](https://github.com/turbot/steampipe-plugin-aws/pull/28))
  ​

_Bug fixes_

- Fixed: Now `bucket_policy_is_public` column for `aws_s3_bucket` will display the correct status of bucket policy ([#36](https://github.com/turbot/steampipe-plugin-aws/pull/36))

_Notes_

- The `lifecycle_rules` column of the table `aws_s3_bucket` has been updated to return an array of lifecycle rules instead of a object with key `Rules` holding lifecycle rules ([#29](https://github.com/turbot/steampipe-plugin-aws/pull/29))

## v0.2.0 [2021-01-28]

_What's new?_
​

- Added: `aws_ssm_parameter` table
  ​
- Updated: `aws_ec2_autoscaling_group` to have `policies` field which contains the details of scaling policy.
- Updated: `aws_ec2_instance` table. Added `instance_status` field which includes status checks, scheduled events and instance state information.
  ​

_Bug fixes_
​

- Fixed: `aws_s3_bucket` table to list buckets even if the region is not set.
