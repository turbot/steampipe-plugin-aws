select transit_gateway_attachment_id, transit_gateway_id, transit_gateway_owner_id, resource_id, resource_type, resource_owner_id, association_state, tags_raw
from aws.aws_ec2_transit_gateway_vpc_attachment
where transit_gateway_attachment_id = '{{ output.resource_id.value }}'
