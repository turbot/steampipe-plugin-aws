select service_name, service_id, service_type, owner, acceptance_required, manages_vpc_endpoints, vpc_endpoint_policy_supported, tags
from aws.aws_vpc_endpoint_service
where service_name = '{{ output.service_name.value }}'
