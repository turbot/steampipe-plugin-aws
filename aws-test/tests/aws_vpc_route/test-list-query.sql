select route_table_id, title
from aws.aws_vpc_route
where title = '{{output.turbot_title.value}}'