select subnet_id, subnet_arn, vpc_id, cidr_block, default_for_az, owner_id, map_customer_owned_ip_on_launch, map_public_ip_on_launch, tags_raw
from aws.aws_vpc_subnet
where subnet_id = '{{ output.resource_id.value }}'
