select arn, name, role_arn, description, event_bus_name,event_pattern,schedule_expression
from aws.aws_cloudwatch_event_rule
where arn = '{{ output.resource_aka.value }}'
