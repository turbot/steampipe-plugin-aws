select network_acl_id, akas, tags, title
from aws.aws_vpc_network_acl
where network_acl_id = '{{ output.resource_id.value }}'
