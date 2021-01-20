select route_table_id, owner_id
from aws.aws_vpc_route_table
where route_table_id = '{{ output.resource_id.value }}'
