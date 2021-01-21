select title, name, rest_api_id, id
from aws.aws_api_gateway_authorizer
where rest_api_id = '{{output.rest_api_id.value}}'
