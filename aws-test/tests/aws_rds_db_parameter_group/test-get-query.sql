select name, arn, description, db_parameter_group_family
from aws.aws_rds_db_parameter_group
where name = '{{ resourceName }}'