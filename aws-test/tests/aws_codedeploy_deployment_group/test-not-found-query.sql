select
  arn
from
  aws_codedeploy_deployment_group
where
  application_name = '{{ output.app_name.value }}' and deployment_group_id = '{{ output.resource_id.value }}-dummy';
