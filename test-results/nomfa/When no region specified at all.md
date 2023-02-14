When no region specified at all

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
| elasticsearch_admin                                                                  | arn:aws:iam::533793682495:policy/turbot/elasticsearch_admin                                                                        |
| ses_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/ses_admin                                                                                  |
| admin_2                                                                              | arn:aws:iam::533793682495:policy/turbot/admin_2                                                                                    |
| braket_operator                                                                      | arn:aws:iam::533793682495:policy/turbot/braket_operator                                                                            |
| AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     | arn:aws:iam::533793682495:policy/service-role/AWSLambdaBasicExecutionRole-f6394e14-e914-44bf-bf63-1c5665b1aee2                     |
| test1                                                                                | arn:aws:iam::533793682495:policy/test1                                                                                             |
| simpledb_operator                                                                    | arn:aws:iam::533793682495:policy/turbot/simpledb_operator                                                                          |
| turbot_deny                                                                          | arn:aws:iam::533793682495:policy/turbot/turbot_deny                                                                                |
| owner_1                                                                              | arn:aws:iam::533793682495:policy/turbot/owner_1                                                                                    |
| ssm_admin                                                                            | arn:aws:iam::533793682495:policy/turbot/ssm_admin                                                                                  |
| greengrass_admin                                                                     | arn:aws:iam::533793682495:policy/turbot/greengrass_admin                                                                           |
| codestar_metadata                                                                    | arn:aws:iam::533793682495:policy/turbot/codestar_metadata                                                                          |
| CodeBuildBasePolicy-test1-us-east-1                                                  | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-test1-us-east-1                                                  |
| rds_owner                                                                            | arn:aws:iam::533793682495:policy/turbot/rds_owner                                                                                  |
| CodeBuildBasePolicy-testbppp-us-east-2                                               | arn:aws:iam::533793682495:policy/service-role/CodeBuildBasePolicy-testbppp-us-east-2                                               |
| logs_operator                                                                        | arn:aws:iam::533793682495:policy/turbot/logs_operator                                                                              |
| test                                                                                 | arn:aws:iam::533793682495:policy/test                                                                                              |
| start-pipeline-execution-us-east-1-test                                              | arn:aws:iam::533793682495:policy/service-role/start-pipeline-execution-us-east-1-test                                              |
| secretsmanager_admin                                                                 | arn:aws:iam::533793682495:policy/turbot/secretsmanager_admin                                                                       |
| workspaces_operator                                                                  | arn:aws:iam::533793682495:policy/turbot/workspaces_operator                                                                        |
| readonly_1                                                                           | arn:aws:iam::533793682495:policy/turbot/readonly_1                                                                                 |
| cloudtrail_metadata                                                                  | arn:aws:iam::533793682495:policy/turbot/cloudtrail_metadata                                                                        |
| ssm_readonly                                                                         | arn:aws:iam::533793682495:policy/turbot/ssm_readonly                                                                               |
| sqs_metadata                                                                         | arn:aws:iam::533793682495:policy/turbot/sqs_metadata                                                                               |
| test_role_to_user                                                                    | arn:aws:iam::533793682495:policy/test_role_to_user                                                                                 |
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
| sqs_admin           | arn:aws:iam::533793682495:policy/turbot/sqs_admin           |
| apigateway_metadata | arn:aws:iam::533793682495:policy/turbot/apigateway_metadata |
| rds_metadata        | arn:aws:iam::533793682495:policy/turbot/rds_metadata        |
| redshift_operator   | arn:aws:iam::533793682495:policy/turbot/redshift_operator   |
| user                | arn:aws:iam::533793682495:policy/turbot/user                |
| owner_1             | arn:aws:iam::533793682495:policy/turbot/owner_1             |
| apigateway_operator | arn:aws:iam::533793682495:policy/turbot/apigateway_operator |
| ssm_operator        | arn:aws:iam::533793682495:policy/turbot/ssm_operator        |
| rds_owner           | arn:aws:iam::533793682495:policy/turbot/rds_owner           |
| readonly_1          | arn:aws:iam::533793682495:policy/turbot/readonly_1          |
+---------------------+-------------------------------------------------------------+