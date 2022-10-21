select title, akas
from aws_route53_record
where zone_id = '{{ output.zone_id.value }}' and name = 'www.{{ resourceName }}.com.'