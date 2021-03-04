select title, akas, tags, region, account_id
from aws.aws_vpc_flow_log
where flow_log_id = '{{ output.flow_log_id.value }}'::asd'