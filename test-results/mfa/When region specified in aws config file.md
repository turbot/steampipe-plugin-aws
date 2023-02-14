When region specified in aws config file

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
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
| apigateway_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_operator                                                                        |
| AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     |
| user                                                                                 | arn:aws:iam::533793682495:policy/turbot/user                                                                                       |
| elasticsearch_operator                                                               | arn:aws:iam::533793682495:policy/turbot/elasticsearch_operator                                                                     |
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| turbot_lockdown_1                                                                    | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_1                                                                          |
| CodeBuildS3ReadOnlyPolicy-codebuild-m-us-east-1                                      | arn:aws:iam::533793682495:policy/service-role/CodeBuildS3ReadOnlyPolicy-codebuild-m-us-east-1                                      |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| lambda-rotate-key-policy                                                             | arn:aws:iam::533793682495:policy/lambda-rotate-key-policy                                                                          |
| turbot_lockdown_restrict_non_owners                                                  | arn:aws:iam::533793682495:policy/turbot/turbot_lockdown_restrict_non_owners                                                        |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| turbot_deny                                                                          | arn:aws:iam::533793682495:policy/turbot/turbot_deny                                                                                |
| rds_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_metadata                                                                               |
| secretsmanager_admin                                                                 | arn:aws:iam::533793682495:policy/turbot/secretsmanager_admin                                                                       |
| config_admin                                                                         | arn:aws:iam::533793682495:policy/turbot/config_admin                                                                               |
| CodeBuildBasePolicy-test-project1-ap-south-1                                         | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-project1-ap-south-1                                         |
| ses_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/ses_admin                                                                                  |
| logs_metadata                                                                        | arn:aws:iam::533793682495:policy/turbot/logs_metadata                                                                              |
| example-policy                                                                       | arn:aws:iam::533793682495:policy/example-policy                                                                                    |
| apigateway_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata                                                                        |
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
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
+---------------------+-------------------------------------------------------------+