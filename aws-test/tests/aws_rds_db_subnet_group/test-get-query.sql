select name, arn, description, vpc_id
from aws.aws_rds_db_subnet_group
where name = '{{ resourceName }}'
