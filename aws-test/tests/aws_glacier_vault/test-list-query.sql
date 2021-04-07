select akas, vault_name, title, tags
from aws.aws_glacier_vault
where akas::text = '["{{ output.resource_aka.value }}"]';