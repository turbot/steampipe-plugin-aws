## v0.32.1 [2021-09-13]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v141--2021-07-20) ([#653](https://github.com/turbot/steampipe-plugin-aws/pull/653))

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
