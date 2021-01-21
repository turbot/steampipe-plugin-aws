select name
from aws_api_gateway_stage
where name = '{{ resourceName }}'
