select akas, region
from aws_securityhub_standards_control
where standards_control_arn = '{{ output.standards_control_arn.value }}';
