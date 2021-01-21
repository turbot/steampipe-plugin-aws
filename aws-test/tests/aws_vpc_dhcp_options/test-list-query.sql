select dhcp_options_id, tags_src, tags, title, akas
from aws.aws_vpc_dhcp_options
where akas::text = '["{{output.resource_aka.value}}"]'
