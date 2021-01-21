select internet_gateway_id, owner_id, attachments, tags_src
from aws.aws_vpc_internet_gateway
where internet_gateway_id = '{{ output.resource_id.value }}'
