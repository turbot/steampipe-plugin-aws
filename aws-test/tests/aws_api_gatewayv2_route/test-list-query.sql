select api_id, title, akas
from aws.aws_api_gatewayv2_route
where akas::text = '["{{output.resource_aka.value}}"]'
