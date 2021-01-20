select allocation_id, public_ip, public_ipv4_pool, domain, tags_raw
from aws.aws_vpc_eip
where allocation_id = '{{ output.resource_id.value }}'
