select title, region, account_id
from aws_service_discovery_instance
where id = '{{ output.resource_id.value }}' and region = '{{ output.aws_region.value }}';
