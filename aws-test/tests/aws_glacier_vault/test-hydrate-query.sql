select vault_name, policy, vault_lock_policy_std, akas, tags, title
from aws.aws_glacier_vault
where vault_name = '{{ output.resource_id.value }}';