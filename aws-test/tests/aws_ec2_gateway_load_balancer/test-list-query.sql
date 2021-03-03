select name, arn
from aws_new.aws_ec2_gateway_load_balancer
where name = '{{resourceName}}'
