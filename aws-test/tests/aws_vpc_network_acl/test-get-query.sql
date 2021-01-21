select network_acl_id, vpc_id, owner_id, is_default, entries, tags_src
from aws.aws_vpc_network_acl
where network_acl_id = '{{ output.resource_id.value }}'
