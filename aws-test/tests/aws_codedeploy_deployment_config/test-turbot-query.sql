select
  title,
  account_id
from
  aws_codedeploy_deployment_config
where
  deployment_config_name = '{{ resourceName }}';

