select transit_gateway_id, akas, tags, title
from aws.aws_ec2_transit_gateway
where transit_gateway_id = '{{ output.transit_gateway_id.value }}'
