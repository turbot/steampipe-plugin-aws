select title, akas
from aws_ec2_reserved_instance
where arn = '{{ output.resource_aka.value }}';