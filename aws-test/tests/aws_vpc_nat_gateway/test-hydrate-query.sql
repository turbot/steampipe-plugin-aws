select nat_gateway_id, akas, tags, title
from aws.aws_vpc_nat_gateway
where nat_gateway_id = '{{ output.resource_id.value }}'
