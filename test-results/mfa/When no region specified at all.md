When no region specified at all

Welcome to Steampipe v0.18.5
For more information, type .help
> .cache clear
> .cache
Error: command needs 1 argument - got 0
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
+----+---------------+-------+
| az | instance_type | count |
+----+---------------+-------+
+----+---------------+-------+
> select
  instance_id,
  monitoring_state
from
  aws_ec2_instance
where
  monitoring_state = 'disabled';
+-------------+------------------+
| instance_id | monitoring_state |
+-------------+------------------+
+-------------+------------------+
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
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| rds_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_metadata                                                                               |
| turbot_lockdown_3                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_3                                                                          |
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| CodeBuildBasePolicy-test-project1-ap-south-1                                         | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-project1-ap-south-1                                         |
| redshift_metadata                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_metadata                                                                          |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| AWSLambdaVPCAccessExecutionRole-58e8b1f4-58fc-4fef-b455-c93ff49ec89c                 | arn:aws:iam::533793682495:policy/service-role/AWSLambdaVPCAccessExecutionRole-58e8b1f4-58fc-4fef-b455-c93ff49ec89c                 |
| sqs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/sqs_admin                                                                                  |
| CodeBuildS3ReadOnlyPolicy-codebuild-m-us-east-1                                      | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-codebuild-m-us-east-1                                      |
| config_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/config_admin                                                                               |
| owner                                                                                | arn:aws:iam::533793682495:policy/turbot/owner                                                                                      |
| CodeBuildS3ReadOnlyPolicy-test-bp1-ap-south-1                                        | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-test-bp1-ap-south-1                                        |
| redshift_owner                                                                       | arn:aws:iam::533793682495:policy/turbot/redshift_owner                                                                             |
| CodeBuildBasePolicy-fsbp-test-ap-south-1                                             | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-fsbp-test-ap-south-1                                             |
| dynamodb_admin                                                                       | arn:aws:iam::533793682495:policy/turbot/dynamodb_admin                                                                             |
| AWSCodePipelineServiceRole-ap-south-1-test-pipeline12072021                          | arn:aws:iam::533793682495:policy/service-role/AWSCodePipelineServiceRole-ap-south-1-test-pipeline12072021                          |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| turbot_lockdown_1                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_1                                                                          |
| operator                                                                             | arn:aws:iam::533793682495:policy/turbot/operator                                                                                   |
| AWSCodePipelineServiceRole-us-east-1-test56                                          | arn:aws:iam::533793682495:policy/service-role/AWSCodePipelineServiceRole-us-east-1-test56                                          |
| CodeBuildBasePolicy-testbppp-us-east-2                                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-testbppp-us-east-2                                               |
| CodeBuildBasePolicy-test-bp1-ap-south-1                                              | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-bp1-ap-south-1                                              |
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
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
+---------------------+-------------------------------------------------------------+