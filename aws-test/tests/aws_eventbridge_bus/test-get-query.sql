select arn, name, tags
from aws.aws_eventbridge_bus
where arn = '{{ output.resource_arn.value }}';
