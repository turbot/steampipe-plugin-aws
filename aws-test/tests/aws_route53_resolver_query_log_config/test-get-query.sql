select id, name, destination_arn
from aws_route53_resolver_query_log_config
where id = '{{ output.resource_id.value }}';