select title, tags, akas
from aws.aws_route53_zone
where id = '{{ output.zone_id.value }}';