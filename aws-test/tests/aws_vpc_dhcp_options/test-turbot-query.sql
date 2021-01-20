select title, tags, akas, partition, region, account_id
from aws.aws_vpc_dhcp_options
where dhcp_options_id = '{{output.resource_id.value}}'