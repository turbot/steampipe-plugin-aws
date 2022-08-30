select title, akas
from aws_route53_traffic_policy_instance
where name = '{{ resourceName }}.{{ resourceName }}.com.';