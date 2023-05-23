select id, service_id
from aws_service_discovery_instance
where id = '{{ output.resource_id.value }}' and service_id = '{{ output.service_id.value }}';
