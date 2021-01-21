select customer_gateway_id, type, bgp_asn, ip_address, tags_src
from aws.aws_vpc_customer_gateway
where customer_gateway_id = '{{ output.resource_id.value }}'
