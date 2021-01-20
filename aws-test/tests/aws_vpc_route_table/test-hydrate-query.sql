select route_table_id, akas, tags, title
from aws.aws_vpc_route_table
where route_table_id = '{{ output.resource_id.value }}'
