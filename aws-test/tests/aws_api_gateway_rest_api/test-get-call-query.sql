select title, akas, tags, policy, policy_std
from aws.aws_api_gateway_rest_api
where api_id = '{{output.resource_id.value}}'
