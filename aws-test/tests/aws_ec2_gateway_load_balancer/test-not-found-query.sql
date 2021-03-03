select name, arn, vpc_id
from aws_new.aws_ec2_gateway_load_balancer
where name = 'dummy-{{ resourceName }}'
