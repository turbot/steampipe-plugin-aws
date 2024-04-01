select replication_task_identifier, arn, title, region
from aws_dms_replication_task
where akas::text = '["{{ output.resource_aka.value }}"]';