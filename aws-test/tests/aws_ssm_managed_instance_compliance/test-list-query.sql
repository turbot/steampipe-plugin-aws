select
  id,
  resource_type,
  name
from
  aws.aws_ssm_managed_instance_compliance
where
  id = '{{ output.resource_id.value }}';