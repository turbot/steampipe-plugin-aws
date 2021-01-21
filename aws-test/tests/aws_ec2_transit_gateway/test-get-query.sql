select transit_gateway_id, transit_gateway_arn, owner_id, description, amazon_side_asn, auto_accept_shared_attachments, default_route_table_association, default_route_table_propagation, dns_support, vpn_ecmp_support, tags_src
from aws.aws_ec2_transit_gateway
where transit_gateway_id = '{{ output.transit_gateway_id.value }}'
