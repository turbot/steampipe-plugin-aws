select title, rest_api_id, resource_id, http_method
from aws_api_gateway_method
where resource_id = '{{output.resource_id.value}}' and rest_api_id = '{{output.rest_api_id.value}}' and http_method = 'POST';