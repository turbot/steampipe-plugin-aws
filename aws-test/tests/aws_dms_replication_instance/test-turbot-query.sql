select title, tags, akas
from aws_dms_replication_instance
where arn = '{{ output.resource_aka.value }}';