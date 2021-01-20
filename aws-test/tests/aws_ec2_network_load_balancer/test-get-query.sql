select name, arn, type, state_code, scheme, ip_address_type, vpc_id
from aws.aws_ec2_network_load_balancer
where arn = '{{ output.resource_aka.value }}'
