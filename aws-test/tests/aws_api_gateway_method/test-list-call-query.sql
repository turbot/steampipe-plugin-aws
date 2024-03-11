select rest_api_id, resource_id, title, http_method
from aws_api_gateway_method
where resource_id = '{{output.resource_id.value}}';
