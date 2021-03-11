select name, arn
from aws.aws_ec2_gateway_load_balancer
where arn = '{{ output.resource_aka.value }}';
