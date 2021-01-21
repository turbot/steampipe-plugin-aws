
select key_manager, aliases, key_rotation_enabled, policy, policy_std, tags_src
from aws.aws_kms_key
where id='{{ output.resource_id.value }}'
