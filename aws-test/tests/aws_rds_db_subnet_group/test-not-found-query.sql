select name, arn, description, status
from aws.aws_rds_db_subnet_group
where name = 'dummy-{{ resourceName }}'