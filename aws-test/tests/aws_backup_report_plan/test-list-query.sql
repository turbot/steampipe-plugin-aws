select name, arn
from aws_backup_report_plan
where arn = '{{ output.resource_aka.value }}';