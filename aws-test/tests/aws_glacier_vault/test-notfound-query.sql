select title, akas, region, account_id
from aws.aws_glacier_vault
where vault_name = '{{ output.resource_id.value }}test';