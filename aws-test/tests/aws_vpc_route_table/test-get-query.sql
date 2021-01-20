select route_table_id, vpc_id, owner_id, associations, routes, tags_raw
from aws.aws_vpc_route_table
where route_table_id = '{{ output.resource_id.value }}'
