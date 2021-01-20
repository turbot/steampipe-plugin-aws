select route_table_id
from aws.aws_vpc_route
where route_table_id = '{{output.resource_id.value}}'
limit 1