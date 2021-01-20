select allocation_id, akas, tags, title
from aws.aws_vpc_eip
where allocation_id = '{{ output.resource_id.value }}'
