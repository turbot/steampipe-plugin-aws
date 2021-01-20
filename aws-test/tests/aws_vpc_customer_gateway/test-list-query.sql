select customer_gateway_id, type
from aws.aws_vpc_customer_gateway
where customer_gateway_id = '{{ output.resource_id.value }}'
