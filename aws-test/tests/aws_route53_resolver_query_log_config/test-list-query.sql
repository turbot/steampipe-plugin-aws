select id, name
from aws_route53_resolver_query_log_config
where name = '{{ resourceName }}';