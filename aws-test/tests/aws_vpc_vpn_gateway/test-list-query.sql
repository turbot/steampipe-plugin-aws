select vpn_gateway_id, type, tags, title, akas
from aws.aws_vpc_vpn_gateway
where vpn_gateway_id = '{{ output.resource_id.value }}'
