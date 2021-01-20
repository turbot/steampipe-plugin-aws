select customer_gateway_id, akas, tags, title
from aws.aws_vpc_customer_gateway
where customer_gateway_id = '{{ output.resource_id.value }}'
