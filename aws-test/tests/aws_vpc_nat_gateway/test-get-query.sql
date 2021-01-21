select nat_gateway_id, nat_gateway_addresses, vpc_id, subnet_id, tags_src
from aws.aws_vpc_nat_gateway
where nat_gateway_id = '{{ output.resource_id.value }}'
