select nat_gateway_id, vpc_id, tags, title, akas
from aws.aws_vpc_nat_gateway
where nat_gateway_id = '{{ output.resource_id.value }}'
