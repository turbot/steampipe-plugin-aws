select internet_gateway_id, akas, tags, title
from aws.aws_vpc_internet_gateway
where internet_gateway_id = '{{ output.resource_id.value }}'
