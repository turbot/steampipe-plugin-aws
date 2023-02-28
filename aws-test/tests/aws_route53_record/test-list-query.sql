select name, zone_id, type, records, set_identifier, ttl, weight
from aws_route53_record
where zone_id = '{{ output.zone_id.value }}' and name = 'www.{{ resourceName }}.com.'