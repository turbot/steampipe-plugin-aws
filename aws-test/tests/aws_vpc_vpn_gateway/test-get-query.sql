select vpn_gateway_id, state, type, amazon_side_asn, vpc_attachments, tags_src
from aws.aws_vpc_vpn_gateway
where vpn_gateway_id = '{{ output.resource_id.value }}'
