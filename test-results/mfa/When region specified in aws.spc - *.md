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

Error: operation error EC2: DescribeInstances, https response error StatusCode: 401, RequestID: 1b0662c4-dd62-494a-9ee0-835e2c521b25, api error AuthFailure: AWS was not able to validate the provided access credentials (SQLSTATE HV000)

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

Error: operation error EC2: DescribeInstances, https response error StatusCode: 401, RequestID: 801fa5b2-ab97-4e61-8376-f0ce94326c69, api error AuthFailure: AWS was not able to validate the provided access credentials (SQLSTATE HV000)

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
| sqs_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/sqs_metadata                                                                               |
| secretsmanager_admin                                                                 | arn:aws:iam::533793682495:policy/turbot/secretsmanager_admin                                                                       |
| turbot_lockdown_4                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_4                                                                          |
| admin_2                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_2                                                                                    |
| sqs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/sqs_admin                                                                                  |
| cloudtrail_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_metadata                                                                        |
| rds_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_admin                                                                                  |
| ssm_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/ssm_admin                                                                                  |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| apigateway_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata                                                                        |
| metadata_1                                                                           | arn:aws:iam::533793682495:policy/turbot/metadata_1                                                                                 |
| readonly                                                                             | arn:aws:iam::533793682495:policy/turbot/readonly                                                                                   |
| braket_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/braket_admin                                                                               |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| admin_1                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_1                                                                                    |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| CodeBuildS3ReadOnlyPolicy-test1-us-east-1                                            | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-test1-us-east-1                                            |
| turbot_lockdown_1                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_1                                                                          |
| owner                                                                                | arn:aws:iam::533793682495:policy/turbot/owner                                                                                      |
| CodeBuildBasePolicy-test-project1-ap-south-1                                         | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-project1-ap-south-1                                         |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| cloudtrail_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_operator                                                                        |
| rds_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_operator                                                                               |
| turbot_lockdown_3                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_3                                                                          |
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
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
+---------------------+-------------------------------------------------------------+