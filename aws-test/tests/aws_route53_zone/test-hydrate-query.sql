select name, tags_src
from aws.aws_route53_zone
where id = '{{ output.zone_id.value }}'