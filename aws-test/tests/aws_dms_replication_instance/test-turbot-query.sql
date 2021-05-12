select title, tags, akas
from aws.aws_dms_replication_instance
where replication_instance_arn = '{{ output.resource_aka.value }}';