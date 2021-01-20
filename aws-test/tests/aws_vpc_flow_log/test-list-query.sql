select flow_log_id, traffic_type
from aws.aws_vpc_flow_log
where flow_log_id = '{{ output.resource_id.value }}'
