select name, arn, title, tags, akas, region, account_id
from aws_secrets_manager_secret
where arn = '{{ output.resource_aka.value }}';