select name, id, status, tags, share_status
from aws_route53_resolver_rule
where akas::text = '["{{ output.resource_aka.value }}"]';
