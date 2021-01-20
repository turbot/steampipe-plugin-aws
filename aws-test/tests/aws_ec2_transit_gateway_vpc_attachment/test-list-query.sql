select transit_gateway_attachment_id, transit_gateway_id
from aws.aws_ec2_transit_gateway_vpc_attachment
where transit_gateway_attachment_id = '{{ output.resource_id.value }}'
