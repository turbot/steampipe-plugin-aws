select tags
from aws_lightsail_instance
where name = '{{ output.resource_name.value }}'
