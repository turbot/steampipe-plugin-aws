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
