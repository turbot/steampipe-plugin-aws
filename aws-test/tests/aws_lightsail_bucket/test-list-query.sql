select name, arn, title, akas
from aws_lightsail_bucket
where akas::text = '["{{ output.resource_aka.value }}"]'
