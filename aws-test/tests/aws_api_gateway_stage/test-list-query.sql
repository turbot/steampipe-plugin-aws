select name, rest_api_id, title, tags, akas
from aws.aws_api_gateway_stage
where akas = '["{{output.resource_aka.value}}"]'
