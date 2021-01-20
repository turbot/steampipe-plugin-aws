select title, akas, tags
from aws.aws_api_gatewayv2_api
where api_id = '{{output.resource_id.value}}'