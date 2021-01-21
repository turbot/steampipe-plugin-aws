select title, akas, tags
from aws.aws_api_gatewayv2_stage
where stage_name = '{{resourceName}}'
and api_id = '{{output.api_id.value}}'
