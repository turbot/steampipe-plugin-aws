select service_name, akas, title
from aws.aws_vpc_endpoint_service
where service_name = '{{ output.service_name.value }}'
