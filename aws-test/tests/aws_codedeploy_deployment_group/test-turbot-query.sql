select
  title,
  tags,
  akas,
  account_id
from
  aws_codedeploy_deployment_group
where
  application_name = '{{ output.app_name.value }}';
