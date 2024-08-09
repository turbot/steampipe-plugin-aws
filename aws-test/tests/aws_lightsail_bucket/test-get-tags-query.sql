select tags
from aws_lightsail_bucket
where name = '{{ output.resource_name.value }}'
