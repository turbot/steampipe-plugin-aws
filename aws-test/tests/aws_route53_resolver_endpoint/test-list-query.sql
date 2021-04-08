select name, id, akas, ip_address_count
from aws.aws_route53_resolver_endpoint
where akas = '["{{ output.resource_aka.value }}"]';