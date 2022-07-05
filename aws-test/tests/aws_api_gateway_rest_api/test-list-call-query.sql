select name, api_id, title, tags, akas
from aws.aws_api_gateway_rest_api
where akas = '["{{output.resource_aka.value}}"]';
