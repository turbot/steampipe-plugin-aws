select title, akas, tags
from aws.aws_api_gateway_domain_name
where domain_name = '{{resourceName}}'
