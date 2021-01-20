select transit_gateway_id, transit_gateway_arn
from aws.aws_ec2_transit_gateway
where transit_gateway_id = '{{ output.transit_gateway_id.value }}'
