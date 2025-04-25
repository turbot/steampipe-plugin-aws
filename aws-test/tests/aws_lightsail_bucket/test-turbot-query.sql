select title, akas, region, account_id
from aws_lightsail_bucket
where name = '{{ output.resource_name.value }}'