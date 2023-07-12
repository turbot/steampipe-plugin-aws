select name, arn
from aws_backup_report_plan
where name = '{{ output.resource_name.value }}'