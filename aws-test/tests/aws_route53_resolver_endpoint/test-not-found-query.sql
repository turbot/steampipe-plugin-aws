select title, akas, tags, direction
from aws_route53_resolver_endpoint
where id = '{{ output.resource_id.value }}::asd';