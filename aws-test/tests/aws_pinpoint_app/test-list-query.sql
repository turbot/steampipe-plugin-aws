select id, arn
from aws_pinpoint_app
where arn = '{{ output.resource_aka.value }}';
