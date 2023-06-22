select name, id, arn
from aws_service_discovery_service
where id = '{{ output.resource_id.value }}';
