select name, title, arn
from aws_lightsail_bucket
where name = '{{ resourceName }}'
