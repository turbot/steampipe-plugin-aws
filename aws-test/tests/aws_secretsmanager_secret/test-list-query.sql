select name, arn, region
from aws_secretsmanager_secret
where akas::text = '["{{ output.resource_aka.value }}"]';