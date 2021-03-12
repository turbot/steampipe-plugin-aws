select backup_plan_id, arn
from aws.aws_backup_plan
where backup_plan_id = '{{ output.id.value }}'