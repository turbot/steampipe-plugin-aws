When region specified in aws config file - us-east-2

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
> 
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
| CodeBuildBasePolicy-testbppp-us-east-2                                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-testbppp-us-east-2                                               |
| ecs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/ecs_admin                                                                                  |
| metadata                                                                             | arn:aws:iam::533793682495:policy/turbot/metadata                                                                                   |
| role_test_s3_admin_same_account                                                      | arn:aws:iam::533793682495:policy/role_test_s3_admin_same_account                                                                   |
| AWSCodePipelineServiceRole-us-east-1-pipeline-m                                      | arn:aws:iam::533793682495:policy/service-role/AWSCodePipelineServiceRole-us-east-1-pipeline-m                                      |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| admin_1                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_1                                                                                    |
| sqs_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/sqs_admin                                                                                  |
| tagging_metadata                                                                     | arn:aws:iam::533793682495:policy/turbot/tagging_metadata                                                                           |
| CodeBuildBuildBatchPolicy-test-cb12345-us-east-2-tests                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBuildBatchPolicy-test-cb12345-us-east-2-tests                               |
| rds_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/rds_operator                                                                               |
| greengrass_admin                                                                     | arn:aws:iam::533793682495:policy/turbot/greengrass_admin                                                                           |
| apigateway_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/apigateway_operator                                                                        |
| ssm_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_readonly                                                                               |
| ssm_operator                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_operator                                                                               |
| DAXtoDynamoDBPolicy                                                                  | arn:aws:iam::533793682495:policy/service-role/DAXtoDynamoDBPolicy                                                                  |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| start-pipeline-execution-us-east-1-test                                              | arn:aws:iam::533793682495:policy/service-role/start-pipeline-execution-us-east-1-test                                              |
| CodeBuildBasePolicy-test-project1-ap-south-1                                         | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test-project1-ap-south-1                                         |
| turbot_deny                                                                          | arn:aws:iam::533793682495:policy/turbot/turbot_deny                                                                                |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
| CodeBuildBasePolicy-fsbp-test-ap-south-1                                             | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-fsbp-test-ap-south-1                                             |
| redshift_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/redshift_operator                                                                          |
| simpledb_readonly                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_readonly                                                                          |
| redshift_owner                                                                       | arn:aws:iam::533793682495:policy/turbot/redshift_owner                                                                             |
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
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
+---------------------+-------------------------------------------------------------+