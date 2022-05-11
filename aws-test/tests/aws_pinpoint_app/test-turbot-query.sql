select title, tags, akas
from aws_pinpoint_app
where id = '{{ output.resource_id.value }}';
