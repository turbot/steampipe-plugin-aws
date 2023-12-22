select id, service_id
from aws_service_discovery_instance
where service_id = '{{ output.service_id.value }}';
