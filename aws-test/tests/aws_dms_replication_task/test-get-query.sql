select replication_task_identifier, arn, migration_type, replication_instance_arn, target_endpoint_arn, tags_src, title, tags, akas, partition, region, account_id
from aws_dms_replication_task
where arn = '{{ output.resource_aka.value }}';