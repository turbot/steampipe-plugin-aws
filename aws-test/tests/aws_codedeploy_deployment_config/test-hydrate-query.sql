select
  deployment_config_id,
  deployment_config_name,
  compute_platform
from
  aws_codedeploy_deployment_config
where
  deployment_config_name = '{{ resourceName }}';
