select name, arn, type, state_code, vpc_id
from aws_new.aws_ec2_gateway_load_balancer
where name = '{{ resourceName }}'
