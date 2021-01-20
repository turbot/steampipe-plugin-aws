select stage_name
from aws.aws_api_gatewayv2_stage
where stage_name = '{{resourceName}}'