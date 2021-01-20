select vpn_gateway_id, akas, tags, title
from aws.aws_vpc_vpn_gateway
where vpn_gateway_id = '{{ output.resource_id.value }}'
