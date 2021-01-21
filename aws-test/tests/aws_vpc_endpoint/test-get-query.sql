select vpc_endpoint_id, title, akas, tags, tags_src
from aws.aws_vpc_endpoint
where vpc_endpoint_id = '{{ output.resource_id.value }}'
