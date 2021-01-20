select network_interface_id, owner_id, requester_managed, source_dest_check, interface_type, description
from aws.aws_ec2_network_interface
where network_interface_id = '{{ output.resource_id.value }}'
