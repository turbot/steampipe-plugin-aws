select transit_gateway_route_table_id, transit_gateway_id, default_association_route_table, default_propagation_route_table, tags_src
from aws.aws_ec2_transit_gateway_route_table
where transit_gateway_route_table_id = '{{ output.resource_id.value }}'
