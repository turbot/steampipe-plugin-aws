select flow_log_id, turbot_akas, turbot_tags, turbot_title
from aws.aws_vpc_flow_log
where flow_log_id = '{{ output.resource_id.value }}'
