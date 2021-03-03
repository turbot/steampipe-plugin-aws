select arn, load_balancer_attributes, tags_src
from aws_new.aws_ec2_gateway_load_balancer
where arn = '{{ output.resource_aka.value }}'