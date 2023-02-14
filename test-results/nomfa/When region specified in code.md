When region specified in code

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
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| apigateway_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_operator                                                                        |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| sqs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/sqs_admin                                                                                  |
| ssm_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_metadata                                                                               |
| CodeBuildBasePolicy-test-cb1234-us-east-2                                            | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-cb1234-us-east-2                                            |
| CodeBuildS3ReadOnlyPolicy-test1-us-east-1                                            | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-test1-us-east-1                                            |
| kms_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/kms_metadata                                                                               |
| ecs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/ecs_admin                                                                                  |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| metadata                                                                             | arn:aws:iam::533793682495:policy/turbot/metadata                                                                                   |
| CodeBuildBasePolicy-fsbp-test-ap-south-1                                             | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-fsbp-test-ap-south-1                                             |
| AWSLambdaVPCAccessExecutionRole-58e8b1f4-58fc-4fef-b455-c93ff49ec89c                 | arn:aws:iam::533793682495:policy/service-role/AWSLambdaVPCAccessExecutionRole-58e8b1f4-58fc-4fef-b455-c93ff49ec89c                 |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| simpledb_admin                                                                       | arn:aws:iam::533793682495:policy/turbot/simpledb_admin                                                                             |
| admin                                                                                | arn:aws:iam::533793682495:policy/turbot/admin                                                                                      |
| cloudtrail_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_operator                                                                        |
| operator_2                                                                           | arn:aws:iam::533793682495:policy/turbot/operator_2                                                                                 |
| cloudtrail_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_metadata                                                                        |
| readonly                                                                             | arn:aws:iam::533793682495:policy/turbot/readonly                                                                                   |
| sqs_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/sqs_metadata                                                                               |
| simpledb_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_operator                                                                          |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| braket_operator                                                                      | arn:aws:iam::533793682495:policy/turbot/braket_operator                                                                            |
> select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed
  and path = '/turbot/' limit 10;
+---------------------+-------------------------------------------------------------+
| name                | arn                                                         |
+---------------------+-------------------------------------------------------------+
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
+---------------------+-------------------------------------------------------------+