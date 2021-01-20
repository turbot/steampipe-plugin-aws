select domain_name, tags, title, akas
from aws.aws_api_gatewayv2_domain_name
where akas::text = '["{{output.resource_aka.value}}"]'
