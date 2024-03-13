select name, region
from aws_securityhub_action_target
where name = '{{ output.name.value }}';
