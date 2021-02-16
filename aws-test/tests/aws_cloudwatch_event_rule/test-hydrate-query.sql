select tags
from aws.aws_cloudwatch_event_rule
where arn = '{{ output.resource_aka.value }}'
