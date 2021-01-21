select vpc_id, vpc_endpoint_id, tags_src, tags, title, akas
from aws.aws_vpc_endpoint
where akas::text = '["{{output.resource_aka.value}}"]'
