select name, id, comment, private_zone, resource_record_set_count
from aws_route53_zone
where id = '{{ output.zone_id.value }}'