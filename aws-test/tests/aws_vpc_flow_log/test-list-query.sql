select akas, flow_log_id, traffic_type, tags_src, tags, title
from aws.aws_vpc_flow_log
where akas::text = '["{{ output.resource_aka.value }}"]'
