select name, arn, instance_type, state
from aws_appstream_fleet
where name = '{{ output.resource_name.value }}'
