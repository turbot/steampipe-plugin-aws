select vault_name, policy, akas, tags, title
from aws.aws_glacier_vault
where vault_name = '{{ output.resource_id.value }}';