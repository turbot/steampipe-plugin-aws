select subnet_id, subnet_arn
from aws.aws_vpc_subnet
where subnet_id = '{{ output.resource_id.value }}'
