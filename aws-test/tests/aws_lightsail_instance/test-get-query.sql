select name, blueprint_id, title, arn
from aws_lightsail_instance
where name = '{{ resourceName }}'
