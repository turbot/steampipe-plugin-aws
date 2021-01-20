select stage_name, api_id, title, tags, akas
from aws.aws_api_gatewayv2_stage
where akas::text = '["{{output.resource_aka.value}}"]'
