select name, rest_api_id, title, akas, tags
from aws.aws_api_gateway_stage
where name = '{{ resourceName }}'
  and rest_api_id = '{{ output.rest_api_id.value }}'
