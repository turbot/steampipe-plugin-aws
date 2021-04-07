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
