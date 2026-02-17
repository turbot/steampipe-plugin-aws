select fleet_id, arn, type, excess_capacity_termination_policy, replace_unhealthy_instances, tags ->> 'Name' as name
from aws.aws_ec2_fleet
where tags ->> 'Name' = '{{ resourceName }}'
