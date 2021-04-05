select name, tags_src, query_logging_configs
from aws.aws_route53_zone
where id = '{{ output.zone_id.value }}';