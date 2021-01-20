select vpc_id, cidr_block
from aws.aws_vpc
where vpc_id = '{{ output.resource_id.value }}'
