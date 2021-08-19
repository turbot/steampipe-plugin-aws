select title, akas
from aws.aws_ec2_reserved_instance
where arn = '{{ output.resource_aka.value }}';