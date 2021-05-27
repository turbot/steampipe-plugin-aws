select integration_id, integration_type, akas
from aws.aws_api_gatewayv2_integration
where integration_id = '{{output.resource_id.value}}'
and api_id = '{{output.api_id.value}}'
