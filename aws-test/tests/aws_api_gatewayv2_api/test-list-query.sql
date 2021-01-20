select name, api_id, api_endpoint, tags, title, akas
from aws.aws_api_gatewayv2_api
where akas::text = '["{{output.resource_aka.value}}"]'
