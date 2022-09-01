select name, id, type, version::varchar
from aws_route53_traffic_policy
where name = '{{ resourceName }}'