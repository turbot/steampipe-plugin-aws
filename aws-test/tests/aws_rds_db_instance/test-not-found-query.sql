select db_instance_identifier, arn, status, resource_id
from aws.aws_rds_db_instance
where db_instance_identifier = 'dummy-{{ resourceName }}'
