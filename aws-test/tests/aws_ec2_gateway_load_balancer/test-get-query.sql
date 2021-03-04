select name, arn, type, state_code, vpc_id, account_id, region
from aws.aws_ec2_gateway_load_balancer
where name = '{{ resourceName }}';
