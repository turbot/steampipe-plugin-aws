select name, title, akas, region, account_id
from aws_service_discovery_service
where name = '{{ resourceName }}' and region = '{{ output.aws_region.value }}';
