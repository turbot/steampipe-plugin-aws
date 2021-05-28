select
  instance_id,
  title,
  akas
from
  aws.aws_ssm_managed_instance
where
  instance_id = '{{ output.resource_id.value }}';