select name, arn, description
from aws.aws_rds_db_option_group
where name = 'dummy-{{ resourceName }}'
