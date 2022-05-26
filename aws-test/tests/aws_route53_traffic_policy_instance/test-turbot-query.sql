select title, akas
from aws.aws_route53_traffic_policy_instance
where name = '{{ resourceName }}.{{ resourceName }}.com.';