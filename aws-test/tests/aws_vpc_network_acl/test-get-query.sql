select network_acl_id, vpc_id, owner_id, is_default, entries, tags_raw
from aws.aws_vpc_network_acl
where network_acl_id = '{{ output.resource_id.value }}'
