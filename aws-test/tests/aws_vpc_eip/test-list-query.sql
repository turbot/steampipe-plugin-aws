select allocation_id, public_ip
from aws.aws_vpc_eip
where allocation_id = '{{ output.resource_id.value }}'
