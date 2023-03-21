select domain_name
from aws.aws_api_gatewayv2_domain_name
where domain_name = '{{resourceName}}'
