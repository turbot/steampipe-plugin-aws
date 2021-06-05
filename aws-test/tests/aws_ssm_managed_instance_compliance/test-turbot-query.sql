select
  title,
  akas
from
  aws.aws_ssm_managed_instance_compliance
where
  id = '{{ output.resource_id.value }}';