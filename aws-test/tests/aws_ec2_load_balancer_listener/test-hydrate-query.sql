select title, akas, ssl_policy_detail
from aws.aws_ec2_load_balancer_listener
where arn = '{{ output.resource_aka.value }}';
