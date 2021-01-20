select name, title, tags, akas
from aws.aws_rds_db_option_group
where name = '{{ resourceName }}'