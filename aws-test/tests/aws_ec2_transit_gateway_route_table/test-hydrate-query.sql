select transit_gateway_route_table_id, akas, tags, title
from aws.aws_ec2_transit_gateway_route_table
where transit_gateway_route_table_id = '{{ output.resource_id.value }}'
