select title, account_id, standards_control_arn
from aws_securityhub_standards_control
where standards_control_arn = '{{ output.standards_control_arn.value }}';