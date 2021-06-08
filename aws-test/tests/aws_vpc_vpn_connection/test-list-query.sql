select vpn_connection_id, type, vpn_gateway_id, customer_gateway_id
from aws.aws_vpc_vpn_connection
where title = '{{ output.resource_name.value }}'
