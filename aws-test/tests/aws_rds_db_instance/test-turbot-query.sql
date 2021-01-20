select db_instance_identifier, title, tags, akas
from aws.aws_rds_db_instance
where db_instance_identifier = '{{ resourceName }}'