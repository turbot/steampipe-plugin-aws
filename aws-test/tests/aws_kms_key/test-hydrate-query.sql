select
  key_manager,
  key_rotation_enabled,
  policy,
  policy_std,
  tags_src,
  alias ->> 'AliasArn' as alias_arn,
  alias ->> 'AliasName' as alias_name,
  alias ->> 'TargetKeyId' as alias_target_key_id
from
  aws.aws_kms_key,
  jsonb_array_elements(aliases) as alias
where
  id = '{{ output.resource_id.value }}';