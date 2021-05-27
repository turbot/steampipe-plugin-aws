select
  instance_id,
  arn,
  resource_type,
  computer_name,
  platform_name,
  platform_type,
  platform_version
from
  aws.aws_ssm_managed_instance
where
  arn = '{{ output.arn.value }}';