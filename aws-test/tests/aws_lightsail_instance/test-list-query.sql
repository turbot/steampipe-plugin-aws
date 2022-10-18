select name, arn, title, akas
from aws_lightsail_instance
where akas::text = '["{{ output.resource_aka.value }}"]'
