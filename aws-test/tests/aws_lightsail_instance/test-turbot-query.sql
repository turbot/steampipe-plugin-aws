select title, akas, region, account_id
from aws_lightsail_instance
where name = '{{ output.resource_name.value }}'