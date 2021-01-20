select name, id, title, tags, akas
from aws.aws_api_gateway_usage_plan
where akas::text = '["{{output.resource_aka.value}}"]'
