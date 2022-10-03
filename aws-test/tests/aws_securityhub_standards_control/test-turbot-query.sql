select akas, region
from aws_securityhub_standards_control
where arn = '{{ output.standards_control_arn.value }}';
