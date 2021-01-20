select service_name, turbot_akas, turbot_tags, turbot_title
from aws.aws_vpc_endpoint_service
where service_name = '{{ output.service_name.value }}'
