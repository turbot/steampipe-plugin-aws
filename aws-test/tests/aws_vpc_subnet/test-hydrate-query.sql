select subnet_id, akas, tags, title
from aws.aws_vpc_subnet
where subnet_id = '{{ output.resource_id.value }}'
