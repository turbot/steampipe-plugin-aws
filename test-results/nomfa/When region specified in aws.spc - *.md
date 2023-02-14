When region specified in aws.spc - *

Welcome to Steampipe v0.18.5
For more information, type .help
> .cache clear
> .cache off
> .search_path a_role_aab_without_mfa
> select
  placement_availability_zone as az,
  instance_type,
  count(*)
from
  aws_ec2_instance
group by
  placement_availability_zone,
  instance_type;
+------------+---------------+-------+
| az         | instance_type | count |
+------------+---------------+-------+
| eu-west-1a | t2.micro      | 1     |
| us-east-2a | t2.medium     | 1     |
+------------+---------------+-------+
> select
  instance_id,
  monitoring_state
from
  aws_ec2_instance
where
  monitoring_state = 'disabled';
+---------------------+------------------+
| instance_id         | monitoring_state |
+---------------------+------------------+
| i-079fe88a8a1cf793d | disabled         |
| i-01d374202295d62c3 | disabled         |
+---------------------+------------------+
> select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed;
+--------------------------------------------------------------------------------------+------------------------------------------------------------------------------------------------------------------------------------+
| name                                                                                 | arn                                                                                                                                |
+--------------------------------------------------------------------------------------+------------------------------------------------------------------------------------------------------------------------------------+
| test                                                                                 | arn:aws:iam::533793682495:policy/test                                                                                              |
| secretsmanager_admin                                                                 | arn:aws:iam::533793682495:policy/turbot/secretsmanager_admin                                                                       |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| sqs_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/sqs_metadata                                                                               |
| test-cw                                                                              | arn:aws:iam::533793682495:policy/test-cw                                                                                           |
| CloudTrailPolicyForCloudWatchLogs_48aa00af-7546-42a0-b4e4-5a87716d13f0               | arn:aws:iam::533793682495:policy/service-role/CloudTrailPolicyForCloudWatchLogs_48aa00af-7546-42a0-b4e4-5a87716d13f0               |
| simpledb_admin                                                                       | arn:aws:iam::533793682495:policy/turbot/simpledb_admin                                                                             |
| AWSLambdaBasicExecutionRole-8685db96-4c2a-4738-b7e9-82cbd9c300f8                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-8685db96-4c2a-4738-b7e9-82cbd9c300f8                     |
| apigateway_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata                                                                        |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
| operator_2                                                                           | arn:aws:iam::533793682495:policy/turbot/operator_2                                                                                 |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| servicecatalog_admin                                                                 | arn:aws:iam::533793682495:policy/turbot/servicecatalog_admin                                                                       |
| ssm_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_readonly                                                                               |
| metadata_1                                                                           | arn:aws:iam::533793682495:policy/turbot/metadata_1                                                                                 |
| vpc_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/vpc_operator                                                                               |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| readonly                                                                             | arn:aws:iam::533793682495:policy/turbot/readonly                                                                                   |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| kms_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/kms_admin                                                                                  |
| cloudtrail_admin                                                                     | arn:aws:iam::533793682495:policy/turbot/cloudtrail_admin                                                                           |
| turbot_lockdown                                                                      | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown                                                                            |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| AWSLambdaBasicExecutionRole-a9c50c0f-dbd2-4ae3-9ed2-9777e2d482b9                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-a9c50c0f-dbd2-4ae3-9ed2-9777e2d482b9                     |
> select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed
  and path = '/turbot/' limit 10
+---------------------+-------------------------------------------------------------+
| name                | arn                                                         |
+---------------------+-------------------------------------------------------------+
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
+---------------------+-------------------------------------------------------------+