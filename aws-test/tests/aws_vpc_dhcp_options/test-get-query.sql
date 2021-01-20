select dhcp_options_id, tags_raw
from aws.aws_vpc_dhcp_options
where dhcp_options_id = '{{output.resource_id.value}}'