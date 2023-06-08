select title, akas
from aws_route53_query_log
where hosted_zone_id = '{{ output.zone_id.value }}';