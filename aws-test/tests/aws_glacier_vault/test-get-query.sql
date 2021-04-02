select vault_name, title, tags_src
from aws.aws_glacier_vault
where vault_name = '{{ output.resource_id.value }}';