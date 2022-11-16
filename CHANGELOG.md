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
