select name, arn, title, tags, akas, region, account_id
from aws_secretsmanager_secret
where arn = '{{ output.resource_aka.value }}';