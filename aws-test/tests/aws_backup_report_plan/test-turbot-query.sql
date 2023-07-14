select title, arn
from aws.aws_backup_report_plan
where report_plan_name = '{{ output.resource_name.value }}';
