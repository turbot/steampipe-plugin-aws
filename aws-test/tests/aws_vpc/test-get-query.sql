select vpc_id, cidr_block, is_default, owner_id, tags_raw
from aws.aws_vpc
where vpc_id = '{{ output.resource_id.value }}'
