select transit_gateway_attachment_id, akas, tags, title
from aws.aws_ec2_transit_gateway_vpc_attachment
where transit_gateway_attachment_id = '{{ output.resource_id.value }}'
