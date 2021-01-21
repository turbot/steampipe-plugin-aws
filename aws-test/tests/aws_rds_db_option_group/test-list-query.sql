select name, arn
from aws.aws_rds_db_option_group
where arn = '{{ output.resource_aka.value }}'
