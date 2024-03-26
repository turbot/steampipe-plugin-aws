select title, tags, akas
from aws_dms_replication_task
where arn = '{{ output.resource_aka.value }}';