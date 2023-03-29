select api_id, route_id
from aws.aws_api_gatewayv2_route
where route_id = '{{output.resource_id.value}}'
