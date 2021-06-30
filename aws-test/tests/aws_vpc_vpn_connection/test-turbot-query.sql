select account_id, akas, title, tags
from aws.aws_vpc_vpn_connection
where vpn_connection_id = '{{ output.resource_id.value }}';