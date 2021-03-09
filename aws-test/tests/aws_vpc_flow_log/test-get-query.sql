select akas, flow_log_id, deliver_logs_permission_arn, log_destination, log_destination_type, log_format, log_group_name, traffic_type, tags_src, tags, title
from aws.aws_vpc_flow_log
where flow_log_id = '{{ output.resource_id.value }}';
