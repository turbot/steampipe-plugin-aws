select integration_id
from aws.aws_api_gatewayv2_integration
where api_id = '{{output.api_id.value}}'
