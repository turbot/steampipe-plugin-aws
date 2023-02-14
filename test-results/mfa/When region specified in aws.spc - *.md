When region specified in aws.spc - *

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

Error: operation error EC2: DescribeInstances, https response error StatusCode: 401, RequestID: a2429cf1-ed5f-4b4e-ab48-c77d629e1260, api error AuthFailure: AWS was not able to validate the provided access credentials (SQLSTATE HV000)

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

Error: operation error EC2: DescribeInstances, https response error StatusCode: 401, RequestID: cf5481fa-e19c-4918-b9c2-21ad368fee33, api error AuthFailure: AWS was not able to validate the provided access credentials (SQLSTATE HV000)

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
| redshift_admin                                                                       | arn:aws:iam::533793682495:policy/turbot/redshift_admin                                                                             |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| backup_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/backup_admin                                                                               |
| CodeBuildBasePolicy-test-cb-us-east-2                                                | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-cb-us-east-2                                                |
| rds_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_readonly                                                                               |
| config_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/config_admin                                                                               |
| CodeBuildBasePolicy-testbpppasfsfsssasd-us-east-2                                    | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-testbpppasfsfsssasd-us-east-2                                    |
| metadata_1                                                                           | arn:aws:iam::533793682495:policy/turbot/metadata_1                                                                                 |
| lambda_readonly                                                                      | arn:aws:iam::533793682495:policy/turbot/lambda_readonly                                                                            |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| operator_2                                                                           | arn:aws:iam::533793682495:policy/turbot/operator_2                                                                                 |
| CodeBuildS3ReadOnlyPolicy-codebuild-m-us-east-1                                      | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-codebuild-m-us-east-1                                      |
| AWSLambdaBasicExecutionRole-e3f21c99-3d86-4025-9686-64d12c6f9f59                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-e3f21c99-3d86-4025-9686-64d12c6f9f59                     |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| CodeBuildBasePolicy-test-graph-cp-rk-ap-south-1                                      | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-graph-cp-rk-ap-south-1                                      |
| cloudtrail_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_metadata                                                                        |
| backup_operator                                                                      | arn:aws:iam::533793682495:policy/turbot/backup_operator                                                                            |
| rds_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_operator                                                                               |
| AWSGlueServiceRole-test1                                                             | arn:aws:iam::533793682495:policy/service-role/AWSGlueServiceRole-test1                                                             |
| admin_3                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_3                                                                                    |
| apigateway_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata                                                                        |
| admin_2                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_2                                                                                    |
| logs_metadata                                                                        | arn:aws:iam::533793682495:policy/turbot/logs_metadata                                                                              |
| tagging_metadata                                                                     | arn:aws:iam::533793682495:policy/turbot/tagging_metadata                                                                           |
| CodeBuildPublicBuildPolicyCWLogs-test-bp1-ap-south-1-codebuild-test-bp1-service-role | arn:aws:iam::533793682495:policy/service-role/CodeBuildPublicBuildPolicyCWLogs-test-bp1-ap-south-1-codebuild-test-bp1-service-role |
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
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
+---------------------+-------------------------------------------------------------+