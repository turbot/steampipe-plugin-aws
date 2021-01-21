select title, name, rest_api_id
from aws.aws_api_gateway_authorizer
where id = '{{output.resource_id.value}}'
  and rest_api_id = '{{output.rest_api_id.value}}'
