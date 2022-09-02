select akas, region, account_id
from aws_securityhub_action_target
where arn = '{{ output.arn.value }}';
