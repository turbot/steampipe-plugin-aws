select backup_plan_id, status
from aws.aws_backup_job
where backup_plan_id = '{{ output.id.value }}'