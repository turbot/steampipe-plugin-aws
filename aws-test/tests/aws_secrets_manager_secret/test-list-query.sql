select name, arn, region
from aws_secrets_manager_secret
where akas::text = '["{{ output.resource_aka.value }}"]';