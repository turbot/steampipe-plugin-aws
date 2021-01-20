select name, arn
from aws.aws_ec2_network_load_balancer
where name = '{{resourceName}}'
