select
  title,
  akas
from
  aws.aws_ssm_managed_instance_compliance
where
  resource_id = '{{ output.resource_id.value }}';