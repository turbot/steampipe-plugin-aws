When region specified in aws config file - us-east-2

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
| owner                                                                                | arn:aws:iam::533793682495:policy/turbot/owner                                                                                      |
| budgets_admin                                                                        | arn:aws:iam::533793682495:policy/turbot/budgets_admin                                                                              |
| config_operator                                                                      | arn:aws:iam::533793682495:policy/turbot/config_operator                                                                            |
| AWSLambdaBasicExecutionRole-ae8e4a57-9896-411c-89ff-c2e9120fb06a                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-ae8e4a57-9896-411c-89ff-c2e9120fb06a                     |
| CodeBuildBasePolicy-test-bpp-us-east-2                                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-bpp-us-east-2                                               |
| acm_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/acm_readonly                                                                               |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| guardduty_admin                                                                      | arn:aws:iam::533793682495:policy/turbot/guardduty_admin                                                                            |
| sns_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/sns_metadata                                                                               |
| logs_readonly                                                                        | arn:aws:iam::533793682495:policy/turbot/logs_readonly                                                                              |
| cloudwatch_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudwatch_operator                                                                        |
| glue_operator                                                                        | arn:aws:iam::533793682495:policy/turbot/glue_operator                                                                              |
| macie_metadata                                                                       | arn:aws:iam::533793682495:policy/turbot/macie_metadata                                                                             |
| ecr_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ecr_operator                                                                               |
| CodeBuildBasePolicy-testbp12-us-east-2                                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-testbp12-us-east-2                                               |
| tagging_operator                                                                     | arn:aws:iam::533793682495:policy/turbot/tagging_operator                                                                           |
| AWSLambdaBasicExecutionRole-00693508-38cf-4301-8d72-c5b8612241c1                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-00693508-38cf-4301-8d72-c5b8612241c1                     |
| CloudTrailPolicyForCloudWatchLogs_9550b270-6e16-4451-bc3b-20b1d2eb21af               | arn:aws:iam::533793682495:policy/service-role/CloudTrailPolicyForCloudWatchLogs_9550b270-6e16-4451-bc3b-20b1d2eb21af               |
| AWSLambdaVPCAccessExecutionRole-7e15757f-9878-4a78-af57-543e86a12e46                 | arn:aws:iam::533793682495:policy/service-role/AWSLambdaVPCAccessExecutionRole-7e15757f-9878-4a78-af57-543e86a12e46                 |
| CodeBuildS3ReadOnlyPolicy-pci-codebuild-us-east-1                                    | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-pci-codebuild-us-east-1                                    |
| braket_metadata                                                                      | arn:aws:iam::533793682495:policy/turbot/braket_metadata                                                                            |
| CodeBuildBasePolicy-new-project-with-oAuth-ap-south-1                                | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-new-project-with-oAuth-ap-south-1                                |
| ecr_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/ecr_owner                                                                                  |
| AWSLambdaS3ExecutionRole-68ba5eda-3c4c-4e29-9465-6f4e2b666fc8                        | arn:aws:iam::533793682495:policy/service-role/AWSLambdaS3ExecutionRole-68ba5eda-3c4c-4e29-9465-6f4e2b666fc8                        |
| CodeBuildCloudWatchLogsPolicy-test-graph-cb-project-us-east-1                        | arn:aws:iam::533793682495:policy/service-role/CodeBuildCloudWatchLogsPolicy-test-graph-cb-project-us-east-1                        |
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
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
+---------------------+-------------------------------------------------------------+