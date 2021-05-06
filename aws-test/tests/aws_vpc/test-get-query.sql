select vpc_id, arn, cidr_block, is_default, owner_id, tags_src
from aws.aws_vpc
where vpc_id = '{{ output.resource_id.value }}';