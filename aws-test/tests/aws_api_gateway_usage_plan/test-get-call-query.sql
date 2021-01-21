select title, akas, tags
from aws_api_gateway_usage_plan
where id = '{{output.resource_id.value}}'
