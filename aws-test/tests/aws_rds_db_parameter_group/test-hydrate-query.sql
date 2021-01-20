select name, parameters, tag_list
from aws.aws_rds_db_parameter_group
where name = '{{ resourceName }}'