select tags, title
from aws.aws_vpc_peering_connection
where id = '{{ output.id.value }}';
