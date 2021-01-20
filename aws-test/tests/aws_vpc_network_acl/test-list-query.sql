select network_acl_id, owner_id, tags, title, akas
from aws.aws_vpc_network_acl
where network_acl_id = '{{ output.resource_id.value }}'
