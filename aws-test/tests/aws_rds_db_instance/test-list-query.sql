select db_instance_identifier, arn, resource_id
from aws.aws_rds_db_instance
where arn = '{{ output.resource_aka.value }}'