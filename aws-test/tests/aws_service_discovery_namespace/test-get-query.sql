select name, id, arn
from aws_service_discovery_namespace
where id = '{{ output.resource_id.value }}';
