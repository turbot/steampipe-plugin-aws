select integration_id, integration_type, api_id, title, akas
from aws.aws_api_gatewayv2_integration
where akas::text = '["{{output.resource_aka.value}}"]'
