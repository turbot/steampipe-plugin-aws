select name, id, type, version
from aws_route53_traffic_policy
where name = 'dummy-{{ resourceName }}'