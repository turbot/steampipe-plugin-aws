select name, arn
from aws.aws_rds_db_subnet_group
where arn = '{{ output.resource_aka.value }}'
