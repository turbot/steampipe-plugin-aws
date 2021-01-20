select internet_gateway_id, owner_id, tags, title
from aws.aws_vpc_internet_gateway
where internet_gateway_id = '{{ output.resource_id.value }}'
