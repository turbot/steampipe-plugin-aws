select title, akas
from aws_route53_traffic_policy
where name = '{{ resourceName }}'