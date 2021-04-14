select backup_plan_id, arn
from aws.aws_backup_plan
where akas::text = '["{{ output.resource_aka.value }}"]';