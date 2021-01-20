select name, id, tags, title, akas
from aws.aws_cloudformation_stack
where akas::text = '["{{ output.resource_aka.value }}"]'