select id, arn, name
from aws_pinpoint_app
where id = '{{ output.resource_id.value }}';
