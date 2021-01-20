select name, arn, description, allows_vpc_and_non_vpc_instance_memberships, engine_name, major_engine_version, options
from aws.aws_rds_db_option_group
where name = '{{ resourceName }}'