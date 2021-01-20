select name, id, created_date, title, akas 
from aws.aws_api_gateway_api_key
where akas = '["{{output.resource_aka.value}}"]'
