select arn, load_balancer_attributes, tags_src
from aws.aws_ec2_network_load_balancer
where arn = '{{ output.resource_aka.value }}'
