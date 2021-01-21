select group_name, group_id, description, vpc_id, owner_id, tags_src
from aws.aws_vpc_security_group
where group_id = '{{ output.resource_id.value }}'
