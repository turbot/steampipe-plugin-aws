select vpc_id, akas, tags, title
from aws.aws_vpc
where vpc_id = '{{ output.resource_id.value }}'
