select replication_instance_identifier, replication_instance_arn, tags_src, title, tags, region
from aws.aws_dms_replication_instance
where replication_instance_arn = '{{ output.resource_aka.value }}';