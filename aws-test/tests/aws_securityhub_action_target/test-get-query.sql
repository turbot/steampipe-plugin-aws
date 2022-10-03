select name, region
from aws_securityhub_action_target
where arn = '{{ output.arn.value }}';
