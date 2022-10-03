select replication_instance_identifier, arn, tags_src, title, tags, region
from aws_dms_replication_instance
where arn = '{{ output.resource_aka.value }}';