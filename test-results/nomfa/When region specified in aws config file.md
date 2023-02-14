When region specified in aws config file

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
| rds_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_metadata                                                                               |
| user                                                                                 | arn:aws:iam::533793682495:policy/turbot/user                                                                                       |
| simpledb_metadata                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_metadata                                                                          |
| CodeBuildS3ReadOnlyPolicy-test1-us-east-1                                            | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-test1-us-east-1                                            |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| CodeBuildBasePolicy-test-project1-ap-south-1                                         | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-project1-ap-south-1                                         |
| sqs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/sqs_admin                                                                                  |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| AWSLambdaVPCAccessExecutionRole-9955b235-4264-4732-94b7-e91c943b769d                 | arn:aws:iam::533793682495:policy/service-role/AWSLambdaVPCAccessExecutionRole-9955b235-4264-4732-94b7-e91c943b769d                 |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| apigateway_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata                                                                        |
| braket_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/braket_admin                                                                               |
| redshift_owner                                                                       | arn:aws:iam::533793682495:policy/turbot/redshift_owner                                                                             |
| rds_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_operator                                                                               |
| apigateway_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_operator                                                                        |
| turbot_lockdown_restrict_non_owners                                                  | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_restrict_non_owners                                                        |
| cloudtrail_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_operator                                                                        |
| simpledb_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_operator                                                                          |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
| AWSLambdaBasicExecutionRole-26c33b0f-551d-48bb-bbaa-abcbed7b36cc                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-26c33b0f-551d-48bb-bbaa-abcbed7b36cc                     |
| AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     |
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
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
+---------------------+-------------------------------------------------------------+