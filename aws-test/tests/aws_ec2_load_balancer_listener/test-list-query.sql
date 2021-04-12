select arn, load_balancer_arn, port, protocol
from aws.aws_ec2_load_balancer_listener
where arn = '{{ output.resource_aka.value }}';
