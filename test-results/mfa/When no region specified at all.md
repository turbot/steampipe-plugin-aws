When no region specified at all

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

Error: you must specify a region in "regions" in ~/.steampipe/config/aws.spc. Edit your connection configuration file and then restart Steampipe. (SQLSTATE HV000)

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

Error: you must specify a region in "regions" in ~/.steampipe/config/aws.spc. Edit your connection configuration file and then restart Steampipe. (SQLSTATE HV000)

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
| start-pipeline-execution-us-east-1-test                                              | arn:aws:iam::533793682495:policy/service-role/start-pipeline-execution-us-east-1-test                                              |
| admin_1                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_1                                                                                    |
| ssm_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_metadata                                                                               |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
| secretsmanager_operator                                                              | arn:aws:iam::533793682495:policy/turbot/secretsmanager_operator                                                                    |
| CloudTrailPolicyForCloudWatchLogs_3c72f47e-e34f-4c8a-a046-70af8243ae10               | arn:aws:iam::533793682495:policy/service-role/CloudTrailPolicyForCloudWatchLogs_3c72f47e-e34f-4c8a-a046-70af8243ae10               |
| CodeBuildCachePolicy-codebuild-m-us-east-1                                           | arn:aws:iam::533793682495:policy/service-role/CodeBuildCachePolicy-codebuild-m-us-east-1                                           |
| guardduty_operator                                                                   | arn:aws:iam::533793682495:policy/turbot/guardduty_operator                                                                         |
| apigateway_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_operator                                                                        |
| AWSCodePipelineServiceRole-ap-south-1-test-pipeline12072021                          | arn:aws:iam::533793682495:policy/service-role/AWSCodePipelineServiceRole-ap-south-1-test-pipeline12072021                          |
| kms_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/kms_metadata                                                                               |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| AWSCodePipelineServiceRole-us-east-1-pipeline-m                                      | arn:aws:iam::533793682495:policy/service-role/AWSCodePipelineServiceRole-us-east-1-pipeline-m                                      |
| CodeBuildBasePolicy-testbppp-us-east-2                                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-testbppp-us-east-2                                               |
| turbot_lockdown_4                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_4                                                                          |
| AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     |
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| user                                                                                 | arn:aws:iam::533793682495:policy/turbot/user                                                                                       |
| rds_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_admin                                                                                  |
| config_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/config_admin                                                                               |
| glue_cat_db_pol                                                                      | arn:aws:iam::533793682495:policy/glue_cat_db_pol                                                                                   |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| simpledb_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_operator                                                                          |
| CodeBuildS3ReadOnlyPolicy-test-bp1-ap-south-1                                        | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-test-bp1-ap-south-1                                        |
| metadata_1                                                                           | arn:aws:iam::533793682495:policy/turbot/metadata_1                                                                                 |
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
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
+---------------------+-------------------------------------------------------------+