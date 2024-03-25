select backup_plan_id, status
from aws_backup_job
where akas::text = '["{{ output.resource_aka.value }}"]';
