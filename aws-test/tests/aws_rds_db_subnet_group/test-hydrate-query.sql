select name, tag_list
from aws.aws_rds_db_subnet_group
where name = '{{ resourceName }}'
