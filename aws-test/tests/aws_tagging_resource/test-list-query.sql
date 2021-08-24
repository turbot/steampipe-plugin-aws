select name, arn
from aws_tagging_resource
where akas::text = '["{{ output.resource_aka.value }}"]';