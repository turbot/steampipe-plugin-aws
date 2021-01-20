select network_interface_id, akas, tags, title
from aws.aws_ec2_network_interface
where network_interface_id = '{{ output.resource_id.value }}'
