select service_name, service_id
from aws.aws_vpc_endpoint_service
where service_name = '{{ output.service_name.value }}'
