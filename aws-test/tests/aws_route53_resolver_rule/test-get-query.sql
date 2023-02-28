select name, tags, title, akas
from aws_route53_resolver_rule
where id = '{{ output.resource_id.value }}';
