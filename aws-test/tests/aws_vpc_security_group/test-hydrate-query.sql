select group_name, akas, tags, title
from aws.aws_vpc_security_group
where group_id = '{{ output.resource_id.value }}'
