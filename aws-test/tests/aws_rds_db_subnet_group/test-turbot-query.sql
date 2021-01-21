select name, title, tags, akas
from aws.aws_rds_db_subnet_group
where name = '{{ resourceName }}'
