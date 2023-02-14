When region specified in code

Welcome to Steampipe v0.18.5
For more information, type .help
> .cache clear
> .cache off
> .search_path a_role_aab_with_mfa
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
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| sqs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/sqs_admin                                                                                  |
| turbot_lockdown_1                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_1                                                                          |
| apigateway_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_operator                                                                        |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| user                                                                                 | arn:aws:iam::533793682495:policy/turbot/user                                                                                       |
| CodeBuildBasePolicy-test1-us-east-1                                                  | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test1-us-east-1                                                  |
| rds_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_metadata                                                                               |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| cloudtrail_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_metadata                                                                        |
| apigateway_admin                                                                     | arn:aws:iam::533793682495:policy/turbot/apigateway_admin                                                                           |
| AWSCodePipelineServiceRole-us-east-1-test56                                          | arn:aws:iam::533793682495:policy/service-role/AWSCodePipelineServiceRole-us-east-1-test56                                          |
| rds_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_admin                                                                                  |
| CodeBuildBasePolicy-test-project1-ap-south-1                                         | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-project1-ap-south-1                                         |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| turbot_lockdown_3                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_3                                                                          |
| simpledb_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_operator                                                                          |
| CodeBuildBasePolicy-fsbp-test-ap-south-1                                             | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-fsbp-test-ap-south-1                                             |
| redshift_metadata                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_metadata                                                                          |
| ssm_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_readonly                                                                               |
| superuser                                                                            | arn:aws:iam::533793682495:policy/turbot/superuser                                                                                  |
| apigateway_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata                                                                        |
| sqs_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/sqs_readonly                                                                               |
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
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
+---------------------+-------------------------------------------------------------+