select group_name, group_id
from aws.aws_vpc_security_group
where group_id = '{{ output.resource_id.value }}'
