select *
from aws_ec2_reserved_instance
where arn = '{{ output.resource_aka.value }}1p000';