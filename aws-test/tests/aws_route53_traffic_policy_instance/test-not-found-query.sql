select name, id, ttl::text, traffic_policy_version::text, traffic_policy_id, hosted_zone_id
from aws_route53_traffic_policy_instance
where name = 'dummy-{{ resourceName }}';