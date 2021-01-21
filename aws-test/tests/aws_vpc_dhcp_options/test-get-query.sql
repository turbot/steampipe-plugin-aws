select dhcp_options_id, tags_src
from aws.aws_vpc_dhcp_options
where dhcp_options_id = '{{output.resource_id.value}}'
