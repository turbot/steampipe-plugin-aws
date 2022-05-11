select title, akas
from aws.aws_route53_traffic_policy
where name = '{{ resourceName }}'