select title
from aws.aws_vpc_route
where route_table_id = '{{output.resource_id.value}}'
order by title desc
limit 1