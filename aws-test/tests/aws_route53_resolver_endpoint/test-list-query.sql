select name, id, akas, ip_address_count
from aws_route53_resolver_endpoint
where akas = '["{{ output.resource_aka.value }}"]';