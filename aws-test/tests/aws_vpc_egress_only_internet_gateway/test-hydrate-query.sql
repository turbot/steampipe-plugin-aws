select id, akas, tags, title
from aws.aws_vpc_egress_only_internet_gateway
where id = '{{ output.resource_id.value }}'
