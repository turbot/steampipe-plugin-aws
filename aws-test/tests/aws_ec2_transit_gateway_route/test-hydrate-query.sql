select transit_gateway_route_table_id, akas, title
from aws.aws_ec2_transit_gateway_route
where transit_gateway_route_table_id = '{{ output.transit_gateway_rtb_id.value }}';
