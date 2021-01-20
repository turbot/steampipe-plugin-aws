select name, scheme, vpc_id, subnets
from aws.aws_ec2_classic_load_balancer
where akas::text = '["{{ output.resource_aka.value }}"]'