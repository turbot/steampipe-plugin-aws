select name, id
from aws_service_discovery_service
where name = '{{ resourceName }}';
