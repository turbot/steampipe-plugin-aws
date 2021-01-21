select title, tags, akas
from aws.aws_ec2_application_load_balancer
where arn = '{{ output.resource_aka.value }}'
